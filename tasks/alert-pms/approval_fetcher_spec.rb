require 'rspec'
require 'json'
require 'webmock/rspec'
require_relative './approval_fetcher.rb'

describe 'ApprovalFetcher' do
  describe '#fetch' do
    let(:issue_url) { 'https://api.github.com/repos/cloudfoundry/cf-final-release-election/issues/43' }

    let(:issues_list_url) { 'https://api.github.com/repos/cloudfoundry/cf-final-release-election/issues' }
    let(:comments_list_url) { "#{issue_url}/comments" }

    let(:issue_request) do
      stub_request(:get, issues_list_url).with(query: {'access_token' => 'GH_ACCESS_TOKEN'})
    end

    let(:comments_request) do
      stub_request(:get, comments_list_url) .with(query: {'access_token' => 'GH_ACCESS_TOKEN'})
    end

    let(:issue_body) { 'This is the text of the issue' }
    let(:issues_response) do
      [
        {html_url: issue_url, comments_url: comments_list_url, body: issue_body}
      ].to_json
    end

    let(:comments_response) { '[]' }

    before do
      issue_request.to_return(body: issues_response)
      comments_request.to_return(body: comments_response)
    end

    it 'returns the url and text of the latest election Github issue' do
      approval_fetcher = ApprovalFetcher.new(access_token: 'GH_ACCESS_TOKEN')

      expect(approval_fetcher.fetch).to eq([issue_url, issue_body])
      expect(issue_request).to have_been_requested
    end

    context 'when a second RC has been proposed in a comment' do
      let(:comment_body) { 'This is the text of the comment' }
      let(:comments_response) do
        [
          {html_url: "#{issue_url}#issuecomment-1", body: "Wrong body"},
          {html_url: "#{issue_url}#issuecomment-2", body: comment_body}
        ].to_json
      end

      it 'returns the url and text for the latest comment' do
        approval_fetcher = ApprovalFetcher.new(access_token: 'GH_ACCESS_TOKEN')

        expect(approval_fetcher.fetch).to eq(["#{issue_url}#issuecomment-2", comment_body])
        expect(issue_request).to have_been_requested
        expect(comments_request).to have_been_requested
      end
    end
  end
end
