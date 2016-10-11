#!/usr/bin/env ruby
require 'yaml'

manifest_filename = './manifest.yml'
output_filename   = './updated_manifest.yml'
env_repo_path, filename = ARGV[0], ARGV[1]
director_ssl_cert_filename = File.expand_path File.join(env_repo_path, "/certs/director-#{filename}.crt")
director_ssl_key_filename = File.expand_path File.join(env_repo_path, "/certs/director-#{filename}.key")
director_cert = File.read(director_ssl_cert_filename)
director_key = File.read(director_ssl_key_filename)
properties = YAML.load_file(manifest_filename)

properties['jobs'].select do |job|
  if job['name'] == 'bosh'
    job['properties']['director']['ssl']['cert'] = director_cert
    job['properties']['director']['ssl']['key'] = director_key
    job['properties']['director']['user_management']['local']['users'] = [
      { 'hm' => 'hm-password'},
      { ENV['DIRECTOR_USERNAME'] => ENV['DIRECTOR_PASSWORD'] }
    ]
    job['networks'].find {|n| n['name'] == 'vip'}['static_ips'] = [ENV['DIRECTOR_IP']]
  end
end

mbus_url = "https://mbus:mbus-password@#{ENV['DIRECTOR_IP']}:6868"
properties['cloud_provider']['ssh_tunnel']['host'] = ENV['DIRECTOR_IP']
properties['cloud_provider']['ssh_tunnel']['private_key'] = File.join(env_repo_path, ENV['DIRECTOR_SSH_KEY_PATH'])
properties['cloud_provider']['mbus'] = mbus_url
properties['cloud_provider']['properties']['agent']['mbus'] = mbus_url

File.open(output_filename, 'w') do |file|
  YAML.dump(properties, file)
end

