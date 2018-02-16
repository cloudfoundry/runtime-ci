#!/usr/bin/env ruby

require 'hashdiff'
require 'yaml'
require_relative './release_changes.rb'
require_relative './renderer.rb'
require_relative './ops_file_finder.rb'

updates = ReleaseUpdates.load_from_files('cf-deployment.yml')

OpsFileFinder.find_ops_files('cf-deployment-release-candidate').each do |opsfile|
  opsfile_updates = ReleaseUpdates.load_from_files("operations/#{opsfile}", opsfile: true)
  updates.merge!(opsfile_updates)
end

puts Renderer.new.render(release_updates: updates)
