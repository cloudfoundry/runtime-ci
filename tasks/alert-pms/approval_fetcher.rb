class ApprovalFetcher
  def initialize(access_token: nil)
    @access_token = access_token
  end

  def fetch
    response_json = make_request("https://api.github.com/repos/cloudfoundry/cf-final-release-election/issues")
    approval_url = response_json[0].fetch('html_url')
    approval_body = response_json[0].fetch('body')

    comment_list_url = response_json[0].fetch('comments_url')
    comment_list = make_request(comment_list_url)
    unless comment_list.empty?
      approval_url = comment_list.last.fetch('html_url')
      approval_body = comment_list.last.fetch('body')
    end

    [approval_url, approval_body]
  end

  private

  def make_request(url)
    uri = URI("#{url}?access_token=#{@access_token}")
    response_body = Net::HTTP.get(uri)
    JSON.parse(response_body)
  end
end
