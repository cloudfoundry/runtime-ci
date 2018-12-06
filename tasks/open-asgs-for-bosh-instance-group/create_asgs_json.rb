#!/usr/bin/ruby -w
# frozen_string_literal: true

require 'json'
require 'open3'

def write_asgs_json(instance_group_ips, destination_file)
  asgs = []

  instance_group_ips.each do |ip|
    asgs << { protocol: 'tcp', destination: ip, ports: '1-65535' }
  end

  File.open(destination_file, 'w') do |f|
    f.write(asgs.to_json)
  end
end

def get_ips_from_bosh_output(instance_group_name)
  instance_ips = []

  stdout, _, exitcode = Open3.capture3('bosh is --json')

  if exitcode != 0
    puts "'bosh is --json' returned an error: #{stdout}"
    exit(1)
  end

  instances = JSON.parse(stdout)['Tables'][0]['Rows']
  instances.each do |is|
    if is['instance'].include? instance_group_name
      instance_ips << is['ips'].split(/\s/).select { |ip| ip.start_with? '10.' }
    end
  end

  instance_ips
end

instance_name = ARGV[0]
destination_file = ARGV[1]

instance_group_ips = get_ips_from_bosh_output(instance_name)
write_asgs_json(instance_group_ips, destination_file)
