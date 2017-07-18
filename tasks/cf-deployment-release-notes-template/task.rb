#!/usr/bin/env ruby

require 'hashdiff'
require 'yaml'
require_relative './release_changes.rb'
require_relative './renderer.rb'

release_candidate_text = File.read('cf-deployment-release-candidate/cf-deployment.yml')
release_candidate = YAML.load(release_candidate_text)

master_text = File.read('cf-deployment-master/cf-deployment.yml')
master = YAML.load(master_text)

release_updates = ReleaseUpdates.new
changeSet = HashDiff.diff(master['releases'], release_candidate['releases'])
changeSet.each do |change|
  release_updates.load_change(change)
end

puts Renderer.new.render(release_updates: release_updates)
