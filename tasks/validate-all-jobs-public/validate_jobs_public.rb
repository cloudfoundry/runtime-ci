require 'yaml'

Dir.glob("pipelines/*.yml").each do |f|
  job_array = YAML.load_file(f)['jobs'].map{ |j| {:name => j['name'], :public => j['public']} }
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
