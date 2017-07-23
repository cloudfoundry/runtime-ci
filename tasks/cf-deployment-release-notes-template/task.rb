#!/usr/bin/env ruby

require 'hashdiff'
require 'yaml'
require_relative './release_changes.rb'
require_relative './renderer.rb'

release_updates = ReleaseUpdates.load_from_files('cf-deployment.yml')
puts Renderer.new.render(release_updates: release_updates)
