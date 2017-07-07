require_relative "./release_team_collection.rb"
require 'rspec'

describe 'ReleaseTeamCollection' do
  describe '#new' do
    it 'fails if an argument is not a ReleaseTeam' do
      expect {
        ReleaseTeamCollection.new('not a release team')
      }.to raise_error "Not a ReleaseTeam"
    end
  end

  describe '#find_team_by_github_handles' do
    let(:anchor_github) { "@anchor_github" }
    let(:anchor_slack) { "@anchor_slack" }
    let(:pm_github) { "@pm_github" }
    let(:pm_slack) { "@pm_slack" }

    let(:teams) do
      ReleaseTeamCollection.new(
        ReleaseTeam.new(
          name: "Team A",
          anchor_github: anchor_github,
          anchor_slack: anchor_slack,
          pm_github: pm_github,
          pm_slack: pm_slack
        )
      )
    end

    it 'returns the team whose anchor and pm match the provided github handles' do
      team = teams.find_team_by_github_handles(anchor_github, pm_github)
      expect(team).not_to be_nil
      expect(team.name).to eq "Team A"
      expect(team.anchor_github).to eq anchor_github
      expect(team.anchor_slack).to eq anchor_slack
      expect(team.pm_github).to eq pm_github
      expect(team.pm_slack).to eq pm_slack
    end

    context "when one of the fields is empty" do
      let(:anchor_github) { '' }
      let(:anchor_slack) { '' }
      let(:pm_github) { "@nebhale" }
      let(:pm_slack) { "@nebhale" }

      it 'still returns the correct team, based on the non-empty handle' do
        team = teams.find_team_by_github_handles('', '@nebhale')
        expect(team).not_to be_nil
        expect(team.name).to eq 'Team A'
        expect(team.anchor_github).to eq ''
        expect(team.anchor_slack).to eq ''
        expect(team.pm_github).to eq '@nebhale'
        expect(team.pm_slack).to eq '@nebhale'
      end

      context 'and the user requests a nil value' do
        it 'matches the nil against the empty field' do
          team = teams.find_team_by_github_handles(nil, '@nebhale')
          expect(team).not_to be_nil
          expect(team.name).to eq 'Team A'
          expect(team.anchor_github).to eq ''
          expect(team.anchor_slack).to eq ''
          expect(team.pm_github).to eq '@nebhale'
          expect(team.pm_slack).to eq '@nebhale'
        end
      end
    end
  end
end
