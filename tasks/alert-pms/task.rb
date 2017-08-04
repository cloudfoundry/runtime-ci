#!/usr/bin/env ruby

require 'uri'
require 'net/http'
require 'json'
require_relative './release_team_collection.rb'
require_relative './release_team_list.rb'
require_relative './alert_message_writer.rb'
require_relative './approval_fetcher.rb'

def filter_approved_releases(issue_body)
  teams = issue_body.split('----')
  teams.select do |team|
    /\:\-1\:/ =~ team
  end
end

puts 'Finding cf-release-final-election issue url'
approval_fetcher = ApprovalFetcher.new(access_token: ENV.fetch('GH_ACCESS_TOKEN'))
approval_url, approval_body = approval_fetcher.fetch

alert_message_writer = AlertMessageWriter.new(release_teams, approval_url)

waiting_approval = filter_approved_releases(approval_body)
waiting_approval.each do |team_section|
  alert_message_writer.write!(team_section)
end
