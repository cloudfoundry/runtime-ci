#!/usr/bin/env ruby

require 'date'
require 'find'
require 'open3'
require 'optparse'
require 'ostruct'
require 'pathname'
require 'yaml'

class CertChecker
  
  attr_accessor :bad_certs
  def initialize(days_left_threshold, verbose=false)
    @days_left_threshold = days_left_threshold
    @verbose = verbose
  end

  def report_problem(input_file)
    puts "\nProblematic certs in file #{input_file}"
  end

  def report_msg_with_times(msg, date_before, date_after)
    fmt = '%Y-%m-%dT%H:%M:%S'
    puts "#{msg}: cert life: #{date_before.strftime(fmt)} -  #{date_after.strftime(fmt)}"
  end

  def check_file(input_file)
    puts "Checking #{input_file}..." if @verbose
    if File.extname(input_file) == ".yml"
      certs = YAML::load(File.open(input_file))
      @found_certs = {}
      find_certs(certs, [])
      puts "Found #{@found_certs.size} cert(s)" if @found_certs.size > 0 && @verbose
    elsif File.extname(input_file) == ".crt"
      @found_certs = {"certificate": File.read(input_file)}
    end

    status = true
    @found_certs.each do |k, cert|
      #puts "CERTIFICATE START"
      puts k if @verbose
      #puts cert
      #puts "CERTIFICATE END"
      stdin, stdout, stderr = Open3.popen3('openssl x509 -text')

      stdin.puts cert
      stdin.close

      output = stdout.read
      errors = stderr.read

      if errors.size > 0
        puts "Error in cert: #{errors}"
        next
      end
      validity_index = output.index("Validity")
      if !validity_index
        puts "Error in cert: no validity section #(output)"
        next
      end
      date_before = DateTime.parse(output[validity_index..validity_index+120].match(/Not Before\s*:\s*(.*)$/)[1])
      date_after = DateTime.parse(output[validity_index..validity_index+120].match(/Not After\s*:\s*(.*)$/)[1])
      if date_before.to_time > Time.now
        if status
          report_problem(input_file)
          status = false
        end
        msg = "!! Cert #{k} isn't ready yet"
        report_msg_with_times(msg, date_before, date_after)
      end
      # modified julian date - day # since Jan 1 4713 BC
      days_left = date_after.mjd - Date.today.mjd
      if days_left < @days_left_threshold
        if status
          report_problem(input_file)
          status = false
        end
        if days_left < 0
          msg = "!! Cert #{k} has expired"
        elsif days_left == 0
          msg = "!! Cert #{k} will expire in today"
        else
          msg = "!! Cert #{k} will expire in #{days_left} day#{days_left == 1 ? "" : "s"}"
        end
        report_msg_with_times(msg, date_before, date_after)
      elsif @verbose
        puts date_before
        puts date_after
      end
    end
    return status
  end

  private

  def find_certs(value, path)
    if value.is_a?(Array)
      value.each_with_index {|v, i| find_certs(v, path + [i.to_s]) }
    elsif value.is_a?(Hash)
      value.each { |k, v| find_certs(v, path + [k]) }
    elsif value.to_s["BEGIN CERTIFICATE"]
      @found_certs[path.join('.')] = value
    end
  end
  
end

def parse(args)
  options = OpenStruct.new
  options.days_left_threshold = 16
  options.path = Dir.pwd
  options.ignore_paths = []

  opt_parser = OptionParser.new do |opts|
    opts.banner = "Usage: check-certs.rb [options]"

    opts.separator ""
    opts.separator "Specific options:"

    opts.on("-d", "--days-left DAYS", Integer,
            "Number of days left for a cert's life before complaining") do |tolerance|
      options.days_left_threshold = tolerance
    end

    opts.on("-p", "--path PATH", String,
            "Path to dir to test") do |path|
      options.path = path
    end

    # this will be passed as YAML fragment from Concourse tasks
    # use "Object" to allow empty value
    opts.on("-i", "--ignore-paths PATH", Object,
            "Paths to ignore (relative to --path)") do |path|
      options.ignore_paths = path
    end

    opts.on_tail("-h", "--help", "Show this message") do
      puts opts
      exit
    end
  end
  opt_parser.parse!(args)
  options
end

options = parse(ARGV)

overall_exit_code = 0
ok_files = %{ 'director-vars-store.yml': true, 'jumpbox-vars-store.yml': true, 'vars.yml': true }

if options.ignore_paths.length > 0
  ignore_paths = YAML::load(options.ignore_paths)
  puts "Ignoring files in #{ignore_paths.length} path(s): #{ignore_paths}"
else
  ignore_paths = []
end

Find.find(options.path) do |input_file|
  next if FileTest.directory?(input_file)
  next if ignore_paths.include? Pathname.new(input_file).relative_path_from(Pathname.new(options.path)).dirname.to_s
  if ok_files[File.basename(input_file)] || File.extname(input_file) == ".crt"
    checker = CertChecker.new(options.days_left_threshold)
    if !checker.check_file(input_file)
      overall_exit_code = 1
    end
  end
end

exit overall_exit_code
