require 'dogapi'
require 'erb'
require 'json'

api_key=ENV.fetch('DATADOG_API_KEY')
app_key=ENV.fetch('APP_KEY')

diego_template_path = "/Users/pivotal/workspace/datadog-config-oss/diego-datadog-templates/#{ENV.fetch("DIEGO_HEALTH_TEMPLATE_PATH")}"
loggregator_template_path="/Users/pivotal/workspace/datadog-config-oss/#{ENV.fetch("LOGGREGATOR_ALERT_TEMPLATE_PATH")}"

dog = Dogapi::Client.new(api_key, app_key)

dashboards = dog.get_all_screenboards().last.fetch("screenboards")
dashboards.each do |board|
  puts "Deleting #{board.fetch('title')}"
  dog.delete_screenboard(board.fetch('id'))
end

timeboards = dog.get_dashboards().last.fetch("dashes")
timeboards.each do |board|
  puts "Deleting #{board.fetch('title')}"
  dog.delete_dashboard(board.fetch('id'))
end

environment = 'A1'
deployment = "cf-a1"
diego_deployment = "cf-a1-diego"
metron_agent_diego_deployment = diego_deployment
diego_dashboard_config = JSON.parse(ERB.new(File.read(diego_template_path)).result(binding))
puts "Creating dashboard #{diego_dashboard_config["title"]}"
response = dog.create_dashboard(diego_dashboard_config["title"],
                     diego_dashboard_config["description"],
                     diego_dashboard_config["graphs"])
if response.first != "200"
  raise "Failed to create dashboard. API response: #{response}"
end

monitors = dog.get_all_monitors().last
monitors.each do |monitor|
  puts "Deleting #{monitor.fetch('name')}"
  dog.delete_monitor(monitor.fetch('id'))
end

metron_agent_deployment=deployment
monitoringAndMetrics_pagerduty="not production deployment"
loggregator_alert_config = JSON.parse(ERB.new(File.read(loggregator_template_path)).result(binding))

puts "Creating monitor #{loggregator_alert_config["name"]}"
response = dog.monitor("metric alert",
            loggregator_alert_config["query"],
            name: loggregator_alert_config["name"],
            message: loggregator_alert_config["message"],
            options: {
              notify_no_data: loggregator_alert_config["notify_no_data"],
              no_data_timeframe: loggregator_alert_config["no_data_timeframe"],
            })
if response.first != "200"
  raise "Failed to create alert. API response: #{response}"
end
