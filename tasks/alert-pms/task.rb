#!/usr/bin/env ruby

require 'uri'
require 'net/http'
require 'json'
require_relative './release_team_collection.rb'

release_teams = ReleaseTeamCollection.new(
  ReleaseTeam.new(
    name: "release-integration",
    anchor_github: '@staylor14',
    anchor_slack:  '@staylor',
    pm_github:     '@dsabeti',
    pm_slack:      '@dsabeti'
  ),
  ReleaseTeam.new(
    name: "infrastructure",
    anchor_github: '@genevievelesperance',
    anchor_slack:  '@genevieve',
    pm_github:     '@evanfarrar',
    pm_slack:      '@efarrar'
  ),
  ReleaseTeam.new(
    name: "postgres-release",
    anchor_github: '@valeriap',
    anchor_slack:  '@valepert',
    pm_github:     '@smoser-ibm',
    pm_slack:      '@smoser'
  ),
  ReleaseTeam.new(
    name: "routing",
    anchor_github: '@shashwathi',
    anchor_slack:  '@sreddy',
    pm_github:     '@shalako',
    pm_slack:      '@shannon'
  ),
  ReleaseTeam.new(
    name: "uaa",
    anchor_github: '@jhamon',
    anchor_slack:  '@jhamon',
    pm_github:     '@sreetummidi',
    pm_slack:      '@stummidi'
  ),
  ReleaseTeam.new(
    name: "loggregator",
    anchor_github: '@jasonkeene',
    anchor_slack:  '@jasonkeene',
    pm_github:     '@ahevenor',
    pm_slack:      '@hev'
  ),
  ReleaseTeam.new(
    name: "capi",
    anchor_github: '@gerg',
    anchor_slack:  '@gerg',
    pm_github:     '@zrob',
    pm_slack:      '@zach'
  ),
  ReleaseTeam.new(
    name: "buildpacks",
    anchor_github: '@dgodd',
    anchor_slack:  '@dogddard',
    pm_github:     '@sclevine',
    pm_slack:      '@slevine'
  ),
  ReleaseTeam.new(
    name: "diego",
    anchor_github: '@jvshahid',
    anchor_slack:  '@jvshahid',
    pm_github:     '@ematpl',
    pm_slack:      '@emalm'
  ),
  ReleaseTeam.new(
    name: "container-networking",
    anchor_github: '@jaydunk',
    anchor_slack:  '@ldunkelberger',
    pm_github:     '@rusha19',
    pm_slack:      '@uramachandran'
  ),
  ReleaseTeam.new(
    name: "grootfs",
    anchor_github: '@spikymonkey',
    anchor_slack:  '@gareth',
    pm_github:     '@glestaris',
    pm_slack:      '@glestaris'
  ),
  ReleaseTeam.new(
    name: "java-buildpack",
    anchor_github: '',
    anchor_slack:  '',
    pm_github:     '@nebhale',
    pm_slack:      '@nebhale'
  )
)

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
