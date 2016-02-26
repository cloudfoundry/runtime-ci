require 'yaml'

my_dir= ENV['MY_DIR']
pipelines_dir="#{my_dir}/../../../pipelines"

Dir.foreach(pipelines_dir) do |f|
  next if File.directory?(f)
  next if f =~ /cf-release.yml/

  job_array = YAML.load_file("#{pipelines_dir}/#{f}")['jobs'].map{ |j| {:name => j['name'], :public => j['public']} }
  job_array.each do |job|
    RSpec.describe "pipeline jobs" do
      context "for #{f}\##{job[:name]}" do
        it "should have public: true" do
          expect(job[:public]).to be_truthy
        end
      end
    end
  end
end
