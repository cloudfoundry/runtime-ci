require 'dogapi'
require 'erb'
require 'json'

api_key=ENV.fetch('DATADOG_API_KEY')
app_key=ENV.fetch('APP_KEY')

dog = Dogapi::Client.new(api_key, app_key)

dashboards = dog.get_all_screenboards().last.fetch('screenboards')
dashboards.each do |board|
  puts "Deleting #{board.fetch('title')}"
  dog.delete_screenboard(board.fetch('id'))
end

timeboards = dog.get_dashboards().last.fetch('dashes')
timeboards.each do |board|
  puts "Deleting #{board.fetch('title')}"
  dog.delete_dashboard(board.fetch('id'))
end

diego_dashboard_config = DiegoHealthDashboard.new().generate_config()

puts "Creating dashboard #{diego_dashboard_config['title']}"
response = dog.create_dashboard(diego_dashboard_config['title'],
                     diego_dashboard_config['description'],
                     diego_dashboard_config['graphs'])
if response.first != '200'
  raise "Failed to create dashboard. API response: #{response}"
end

monitors = dog.get_all_monitors().last
monitors.each do |monitor|
  puts "Deleting #{monitor.fetch('name')}"
  dog.delete_monitor(monitor.fetch('id'))
end

loggregator_alert_config = LoggregatorMonitor.new().generate_config

puts "Creating monitor #{loggregator_alert_config['name']}"
response = dog.monitor('metric alert',
            loggregator_alert_config['query'],
            name: loggregator_alert_config['name'],
            message: loggregator_alert_config['message'],
            options: {
              notify_no_data: loggregator_alert_config['notify_no_data'],
              no_data_timeframe: loggregator_alert_config['no_data_timeframe'],
            })
if response.first != '200'
  raise "Failed to create dashboard. API response: #{response}"
end

class DiegoHealthDashboard
  def environment
    ENV.fetch("ENVIRONMENT_DISPLAY_NAME")
  end

  def deployment
    ENV.fetch("CF_DEPLOYMENT_NAME")
  end

  def diego_deployment
    ENV.fetch("DIEGO_DEPLOYMENT_NAME")
  end

  def metron_agent_diego_deployment
    ENV.fetch("METRON_AGENT_DIEGO_DEPLOYMENT_TAG")
  end

  def generate_config
    diego_template_path = "datadog-diego-health-template/#{ENV.fetch('DIEGO_HEALTH_TEMPLATE_PATH')}"
    JSON.parse(ERB.new(File.read(diego_template_path)).result(binding))
  end
end

class LoggregatorMonitor
  def environment
    ENV.fetch("ENVIRONMENT_DISPLAY_NAME")
  end

  def metron_agent_diego_deployment
    ENV.fetch("METRON_AGENT_DIEGO_DEPLOYMENT_TAG")
  end

  def metron_agent_deployment
    ENV.fetch("METRON_AGENT_CF_DEPLOYMENT_TAG")
  end

  def monitoringAndMetrics_pagerduty
    'not production deployment'
  end

  def generate_config
    loggregator_template_path = "datadog-loggregator-alert-template/#{ENV.fetch('LOGGREGATOR_ALERT_TEMPLATE_PATH')}"
    JSON.parse(ERB.new(File.read(loggregator_template_path)).result(binding))
  end
end
