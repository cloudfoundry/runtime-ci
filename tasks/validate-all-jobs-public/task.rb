#!/usr/bin/env ruby

require 'yaml'

def find_yaml_files(directory)
  yaml_files = Dir.glob(File.join(directory, '**', '*.yaml'))
  yml_files = Dir.glob(File.join(directory, '**', '*.yml'))
  yaml_files + yml_files
end

directory = ENV['RUNTIME_CI_DIR']
if directory.nil? || directory.strip.empty?
  puts "Error: Please set the RUNTIME_CI_DIR environment variable to specify the directory."
end

yaml_files = find_yaml_files(directory)
puts "Found #{yaml_files.size} YAML files in #{directory}"

failure = false

yaml_files.each do |file|
  puts "Checking #{file}..."
  yaml_content = YAML.load_file(file, aliases: true)
  next unless yaml_content.is_a?(Hash) && yaml_content.key?('jobs')
  yaml_content['jobs'].each do |job|
    unless job['public'] == true
      puts "- #{job['name']} is not public."
      failure = true
    end
  end
end

if failure
  puts 'Error: Some jobs are not public'
  exit 1
end

puts 'All jobs are public!'
