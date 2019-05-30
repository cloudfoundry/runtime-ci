#!/usr/bin/env ruby

require 'hashdiff'
require 'yaml'
require_relative './release_changes.rb'
require_relative './renderer.rb'
require_relative './ops_file_finder.rb'

updates = BinaryUpdates.load_from_file('dockerfiles/cf-deployment-concourse-tasks/Dockerfile')
puts Renderer.new.render(binary_updates: updates)
