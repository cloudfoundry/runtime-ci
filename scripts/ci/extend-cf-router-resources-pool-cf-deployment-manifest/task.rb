#!/usr/bin/env ruby
require 'yaml'
require 'active_support'

input_filename = './generated-deployment-manifest/' + ENV.fetch('CF_DEPLOYMENT_MANIFEST')
output_filename = './extended-cf-deployment-manifest/' + ENV.fetch('EXTENDED_CF_DEPLOYMENT_MANIFEST')

def get_security_groups
  File.open('bbl-security-groups/' + ENV.fetch('BBL_SECURITY_GROUPS_FILE'), "r") do |f|
    f.map do |line|
       line.split(" ")[1]
     end
  end
end

def get_elb(name)
  IO.popen(%(bbl --state-dir ./env-repo/#{ENV.fetch('BBL_STATE_DIRECTORY')} lbs | grep "#{name}" | cut -d: -f2 | cut -d\  -f2), "r") do |cmd|
    return cmd.read.chomp
  end
end

def merge_resource_pool!(lb_name, resource_pool)
  resource_pool.deep_merge!('cloud_properties' => {'security_groups' => get_security_groups })
  resource_pool.deep_merge!('cloud_properties' => {'elbs' => get_elb(lb_name)})
end

properties = YAML.load_file(input_filename)

def add_security_group_to_router(resource_pool)
  if ["router_z1", "router_z2"].include? resource_pool["name"]
    merge_resource_pool!("CF Router LB", resource_pool)
  end
  if ["access_z1", "access_z2"].include? resource_pool["name"]
    merge_resource_pool!("CF SSH Proxy LB", resource_pool)
  end
end

properties.each do | property_key, property_value |
  if property_key == 'resource_pools'
    property_value.each do | resource_pool_value |
      add_security_group_to_router(resource_pool_value)
    end
  end
end

File.open(output_filename, 'w') do |file|
  YAML.dump(properties, file)
end
