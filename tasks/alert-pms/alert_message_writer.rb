class AlertMessageWriter
  def initialize(release_teams, issue_url)
    @release_teams = release_teams
    @issue_url = issue_url
  end

  def write!(team_section)
    pm_github, anchor_github = team_section.split("\r\n")[1].gsub(':', '').split(', ')
    team = @release_teams.find_team_by_github_handles(anchor_github, pm_github)
    open(File.join("slack-messages", team.name), 'w') do |file|
      file.puts build_message(team)
    end
  end

  private

  def build_message(team)
    message = "Could you please take a look at the latest release candidate: #{@issue_url} cc <@dsabeti>"
    greeting = "Hey there"
    addressees = ""
    addressees << "<#{team.pm_slack}>"
    addressees << " <#{team.anchor_slack}>" unless team.anchor_slack.nil? || team.anchor_slack.empty?
    "#{greeting} #{addressees}. #{message}"
  end
end
