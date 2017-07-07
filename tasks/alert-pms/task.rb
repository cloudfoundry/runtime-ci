#!/usr/bin/env ruby

require 'uri'
require 'net/http'
require 'json'
require_relative './release_team_collection.rb'
require_relative './release_team_list.rb'

def filter_approved_releases(issue_body)
  teams = issue_body.split('----')
  teams.select do |team|
    /\:\-1\:/ =~ team
  end
end

def write_alert_message(team_section, release_teams, issue_url)
  pm_github, anchor_github = team_section.split("\r\n")[1].gsub(':', '').split(', ')
  team = release_teams.find_team_by_github_handles(anchor_github, pm_github)
  open(File.join("slack-messages", team.name), 'w') do |file|
    file.puts "Hey there <#{team.pm_slack}> <#{team.anchor_slack}>. Could you please take a look at the latest release candidate: #{issue_url} cc <@dsabeti>"
  end
end

puts "Finding cf-release-final-election issue url"
uri = URI("https://api.github.com/repos/cloudfoundry/cf-final-release-election/issues?access_token=#{ENV.fetch('GH_ACCESS_TOKEN')}")
response_body = Net::HTTP.get(uri)
response_json = JSON.parse(response_body)
issue_url = response_json[0].fetch('html_url')

issue_body = response_json[0].fetch('body')
waiting_approval = filter_approved_releases(issue_body)
waiting_approval.each do |team_section|
  write_alert_message(team_section, release_teams, issue_url)
end
