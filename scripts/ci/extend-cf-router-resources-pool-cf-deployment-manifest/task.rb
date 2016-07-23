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

def get_elbs
  File.open('bbl-elbs/' + ENV.fetch('BBL_ELBS_FILE'), "r") do |f|
    f.map do |line|
      line.split(" ")[1]
    end
  end
end

properties = YAML.load_file(input_filename)

def add_security_group_to_router(resource_pool)
  if ["router_z1", "router_z2"].include? resource_pool["name"]
    resource_pool.deep_merge!('cloud_properties' => {'security_groups' => get_security_groups })
    resource_pool.deep_merge!('cloud_properties' => {'elbs' => get_elbs })
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
