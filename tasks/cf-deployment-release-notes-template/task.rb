#!/usr/bin/env ruby

require 'hashdiff'
require 'yaml'
require_relative './release_changes.rb'
require_relative './renderer.rb'

updates = ReleaseUpdates.load_from_files('cf-deployment.yml')

opsfile_list = Dir.glob(File.join("cf-deployment-release-candidate", "operations", "*"))
opsfile_list.select! { |opsfile| File.file?(opsfile) }
opsfile_list.map! { |opsfile| opsfile.gsub!('cf-deployment-release-candidate/', '') }

opsfile_list.each do |opsfile|
  opsfile_updates = ReleaseUpdates.load_from_files(opsfile, opsfile: true)
  updates.merge!(opsfile_updates)
end

puts Renderer.new.render(release_updates: updates)
