require 'rspec'
require_relative './ops_file_finder.rb'

describe 'OpsFileFinder' do
  subject { OpsFileFinder.find_ops_files('test-cf-deployment') }

  before(:all) do
    @current_work_dir = Dir.pwd
    @tmp_work_dir = Dir.mktmpdir('test-cf-deployment')

    Dir.chdir(@tmp_work_dir)
    FileUtils.mkdir_p('test-cf-deployment/operations')
  end

  after(:all) do
    Dir.chdir(@current_work_dir)
    FileUtils.rm_rf(@tmp_work_dir) if File.exist?(@tmp_work_dir)
  end

  context 'when there are ops-files in the lop-level operations directory' do
    before do
      File.open('test-cf-deployment/operations/ops.yml', 'w')
      File.open('test-cf-deployment/operations/ops2.yml', 'w')
      File.open('test-cf-deployment/operations/README.md', 'w')
    end

    it 'returns the file names of the ops-files without any additional path' do
      expect(subject).to include 'ops.yml'
      expect(subject).to include 'ops2.yml'
    end

    it 'does not return any files that are not yaml files' do
      expect(subject).not_to include 'README.md'
    end
  end

  context 'when there is another directory in the operations directory' do
    before do
      FileUtils.mkdir_p("test-cf-deployment/operations/#{subfolder}")
      File.open("test-cf-deployment/operations/#{subfolder}/ops.yml", 'w')
      File.open("test-cf-deployment/operations/#{subfolder}/ops2.yml", 'w')
    end

    let(:expect_ops_files_to_be_included) do
      expect(subject).to include "#{subfolder}/ops.yml"
      expect(subject).to include "#{subfolder}/ops2.yml"
    end

    let(:expect_ops_files_to_be_excluded) do
      expect(subject).not_to include "#{subfolder}/ops.yml"
      expect(subject).not_to include "#{subfolder}/ops2.yml"
    end

    let(:subfolder) { 'nested' }

    it 'does not include the directory in the list of files' do
      expect(subject).not_to include 'nested'
    end

    context 'and the directory is the experimental directory' do
      let(:subfolder) { 'experimental' }

      it 'returns the ops-files prepended with "experimentlal"' do
        expect_ops_files_to_be_included
      end
    end

    context 'and the directory is the addons directory' do
      let(:subfolder) { 'addons' }

      it 'returns the ops-files prepended with "addons"' do
        expect_ops_files_to_be_included
      end
    end

    context 'and the directory is the legacy directory' do
      let(:subfolder) { 'legacy' }

      it 'does not return the ops-files prepended with "legacy"' do
        expect_ops_files_to_be_excluded
      end
    end

    context 'and the directory is the workaround directory' do
      let(:subfolder) { 'workaround' }

      it 'does not return the ops-files prepended with "workaround"' do
        expect_ops_files_to_be_excluded
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
