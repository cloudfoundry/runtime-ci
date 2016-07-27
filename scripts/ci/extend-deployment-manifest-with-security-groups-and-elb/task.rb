#!/usr/bin/env ruby
require 'yaml'
require 'active_support'

def aws_output_mapping
  @aws_output_mapping ||= File.open('logical-physical-id/' + ENV.fetch('CLOUDFORMATION_MAPPING_FILE'), "r") do |f|
    mapping = {}
    f.map do |line|
      key, value= line.split(" ")
      mapping[key] = value
    end
    mapping
  end
end

def get_elb(name)
  IO.popen(%(bbl --state-dir ./env-repo/#{ENV.fetch('BBL_STATE_DIRECTORY')} lbs | grep "#{name}" | cut -d: -f2 | cut -d\\  -f2), "r") do |cmd|
    return cmd.read.chomp
  end
end

input_filename = './generated-deployment-manifest/' + ENV.fetch('CF_DEPLOYMENT_MANIFEST')
output_filename = './extended-cf-deployment-manifest/' + ENV.fetch('EXTENDED_CF_DEPLOYMENT_MANIFEST')
properties = YAML.load_file(input_filename)

properties.deep_merge!('properties' => {
  'template_only' => {
    'aws' => {
      'access_key_id' => ENV.fetch('AWS_ACCESS_KEY_ID'),
      'secret_access_key' => ENV.fetch('AWS_SECRET_ACCESS_KEY'),
    }
  }
})

access_pools = %w[access_z1 access_z2]
access_resource_pools = properties['resource_pools'].select do |resource_pool|
  access_pools.include?(resource_pool['name'])
end

access_resource_pools.each do |pool|
  security_groups = [
    aws_output_mapping['CFSSHProxyInternalSecurityGroup'],
    aws_output_mapping['InternalSecurityGroup']
  ]
  pool.deep_merge!('cloud_properties' => {'security_groups' => security_groups})
  pool.deep_merge!('cloud_properties' => {'elbs' => [get_elb('CF SSH Proxy LB')]})
end

router_pools = %w[router_z1 router_z2]
router_resource_pools = properties['resource_pools'].select do |resource_pool|
  router_pools.include?(resource_pool['name'])
end

router_resource_pools.each do |pool|
  security_groups = [
    aws_output_mapping['CFRouterInternalSecurityGroup'],
    aws_output_mapping['InternalSecurityGroup']
  ]
  pool.deep_merge!('cloud_properties' => {'security_groups' => security_groups})
  pool.deep_merge!('cloud_properties' => {'elbs' => [get_elb('CF Router LB')]})
end

File.open(output_filename, 'w') do |file|
  YAML.dump(properties, file)
end
