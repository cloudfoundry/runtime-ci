require 'rspec'
require_relative './alert_message_writer.rb'

describe 'AlertMessageWriter' do
  describe '#write!' do
    let(:fake_release_teams) do
      spy('release_teams').tap do |release_teams_spy|
        allow(release_teams_spy).to receive(:find_team_by_github_handles).and_return(team)
      end
    end

    let(:pm_slack) { '@pm' }
    let(:anchor_slack) { '@anchor' }
    let(:team_name) { 'cf-team' }
    let(:team) do
      spy('team').tap do |team|
        allow(team).to receive(:name).and_return team_name
        allow(team).to receive(:pm_slack).and_return pm_slack
        allow(team).to receive(:anchor_slack).and_return anchor_slack
      end
    end

    let(:team_section) { "\r\n@pm, @anchor:\r\nDoes blah" }

    subject(:alert_message_writer) { AlertMessageWriter.new(fake_release_teams, "issue.url") }

    it 'extacts the github handles and fetches the corresponding slack handles' do
      alert_message_writer.write!(team_section)
      expect(fake_release_teams).to have_received(:find_team_by_github_handles).with('@anchor', '@pm')
    end

    it 'writes the message that mentions both people in a file named after the team' do
      alert_message_writer.write!(team_section)

      file_path = File.join('slack-messages', team_name)
      expect(File.exists?(file_path)).to be true
      expect(File.read(file_path)).to include '<@anchor>'
      expect(File.read(file_path)).to include '<@pm>'
    end

    context 'when the anchor is empty' do
      let(:anchor_slack) { '' }

      it 'only mentions the PM in the message' do
        alert_message_writer.write!(team_section)

        file_path = File.join('slack-messages', team_name)
        expect(File.exists?(file_path)).to be true
        expect(File.read(file_path)).not_to include '@anchor'
        expect(File.read(file_path)).not_to include '<>'

        expect(File.read(file_path)).to include '<@pm>'
      end
    end

    context 'when the anchor is nil' do
      let(:anchor_slack) { nil }

      it 'only mentions the PM in the message' do
        alert_message_writer.write!(team_section)

        file_path = File.join('slack-messages', team_name)
        expect(File.exists?(file_path)).to be true
        expect(File.read(file_path)).not_to include '@anchor'
        expect(File.read(file_path)).not_to include '<>'

        expect(File.read(file_path)).to include '<@pm>'
      end
    end
  end
end
