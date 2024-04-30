require 'rspec'
require_relative './ops_file_finder.rb'

describe 'OpsFileFinder' do
  subject { OpsFileFinder.find_ops_files('test-cf-deployment') }

  before(:all) do
    @current_work_dir = Dir.pwd
    @tmp_work_dir = Dir.mktmpdir('test-cf-deployment')

    Dir.chdir(@tmp_work_dir)
    FileUtils.mkdir_p('test-cf-deployment/operations/test')
  end

  after(:all) do
    Dir.chdir(@current_work_dir)
    FileUtils.rm_rf(@tmp_work_dir) if File.exist?(@tmp_work_dir)
  end

  context 'when there are ops-files in the lop-level operations directory' do
    before do
      File.open('test-cf-deployment/operations/ops.yml', 'w')
      File.open('test-cf-deployment/operations/ops2.yml', 'w')
      File.open('test-cf-deployment/operations/use-compiled-releases.yml', 'w')
      File.open('test-cf-deployment/operations/test/fips-stemcell.yml', 'w')
      File.open('test-cf-deployment/operations/README.md', 'w')
    end

    it 'returns the file names of the ops-files without any additional path' do
      expect(subject).to include 'ops.yml'
      expect(subject).to include 'ops2.yml'
    end

    it 'does not return any files that are not yaml files' do
      expect(subject).not_to include 'README.md'
    end

    context 'and the file is explicitly excluded' do
      it 'does not return the use-compiled-releases.yml ops-file' do
        expect(subject).not_to include "use-compiled-releases.yml"
      end
      it 'does not return the fips-stemcell.yml ops-file' do
        expect(subject).not_to include "test/fips-stemcell.yml"
      end
    end
  end

  context 'when there is another directory in the operations directory' do
    before do
      FileUtils.mkdir_p("test-cf-deployment/operations/#{subfolder}")
      File.open("test-cf-deployment/operations/#{subfolder}/ops.yml", 'w')
      File.open("test-cf-deployment/operations/#{subfolder}/ops2.yml", 'w')
    end

    let(:subfolder) { 'nested' }

    it 'does not include the directory in the list of files' do
      expect(subject).not_to include 'nested'
    end

    context 'and the directory is the experimental directory' do
      let(:subfolder) { 'experimental' }

      it 'returns the ops-files prepended with "experimentlal"' do
        expect(subject).to include "#{subfolder}/ops.yml"
        expect(subject).to include "#{subfolder}/ops2.yml"
      end
    end

    context 'and the directory is the addons directory' do
      let(:subfolder) { 'addons' }

      it 'returns the ops-files prepended with "addons"' do
        expect(subject).to include "#{subfolder}/ops.yml"
        expect(subject).to include "#{subfolder}/ops2.yml"
      end
    end

    context 'directory ends with .yml' do
      let(:subfolder) { 'yaml-dir.yml' }

      it 'does not include the directory in a list of files' do
        expect(subject).not_to include 'yaml-dir.yml'
      end
    end
  end
end
