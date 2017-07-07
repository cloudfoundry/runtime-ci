#!/usr/bin/env ruby

require 'uri'
require 'net/http'
require 'json'
require_relative './release_team_collection.rb'
require_relative './release_team_list.rb'
require_relative './alert_message_writer.rb'

def filter_approved_releases(issue_body)
  teams = issue_body.split('----')
  teams.select do |team|
    /\:\-1\:/ =~ team
  end
end

puts "Finding cf-release-final-election issue url"
uri = URI("https://api.github.com/repos/cloudfoundry/cf-final-release-election/issues?access_token=#{ENV.fetch('GH_ACCESS_TOKEN')}")
response_body = Net::HTTP.get(uri)
response_json = JSON.parse(response_body)
issue_url = response_json[0].fetch('html_url')

alert_message_writer = AlertMessageWriter.new(release_teams, issue_url)

issue_body = response_json[0].fetch('body')
waiting_approval = filter_approved_releases(issue_body)
waiting_approval.each do |team_section|
  alert_message_writer.write!(team_section)
end
