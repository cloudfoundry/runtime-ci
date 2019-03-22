require 'rspec'
require_relative './ops_file_finder.rb'

describe 'OpsFileFinder' do
  before(:all) do
    @current_work_dir = Dir.pwd
    @tmp_work_dir = Dir.mktmpdir('test-cf-deployment')

    Dir.chdir(@tmp_work_dir)

    FileUtils.mkdir_p('test-cf-deployment')
    FileUtils.mkdir_p('test-cf-deployment/operations/experimental')
    FileUtils.mkdir_p('test-cf-deployment/operations/addons')
  end

  after(:all) do
    Dir.chdir(@current_work_dir)
    FileUtils.rm_rf(@tmp_work_dir) if File.exists?(@tmp_work_dir)
  end

  context 'when there are ops-files in the lop-level operations directory' do
    before do
      File.open('test-cf-deployment/operations/ops.yml', 'w') {}
      File.open('test-cf-deployment/operations/ops2.yml', 'w') {}
      File.open('test-cf-deployment/operations/README.md', 'w') {}
    end

    it 'returns the file names of the ops-files without any additional path' do
      ops = OpsFileFinder.find_ops_files('test-cf-deployment')

      expect(ops).to include 'ops.yml'
      expect(ops).to include 'ops2.yml'
    end

    it 'does not return any files that are not yaml files' do
      ops = OpsFileFinder.find_ops_files('test-cf-deployment')

      expect(ops).not_to include 'README.md'
    end
  end

  context 'when there is another directory in the operations directory' do
    before do
      File.open('test-cf-deployment/operations/experimental/ops.yml', 'w') {}
      File.open('test-cf-deployment/operations/experimental/ops2.yml', 'w') {}
      File.open('test-cf-deployment/operations/addons/ops.yml', 'w') {}
      File.open('test-cf-deployment/operations/addons/ops2.yml', 'w') {}
    end

    it 'does not include the directory in the list of files' do
      ops = OpsFileFinder.find_ops_files('test-cf-deployment')

      expect(ops).not_to include 'experimental'
      expect(ops).not_to include 'addons'
    end

    context 'and the directory is the experimental directory' do
      it 'returns the ops-files prepended with "experimentlal"' do
        ops = OpsFileFinder.find_ops_files('test-cf-deployment')

        expect(ops).to include 'experimental/ops.yml'
        expect(ops).to include 'experimental/ops2.yml'
      end
    end

    context 'and the directory is the addons directory' do
      it 'returns the ops-files prepended with "addons"' do
        ops = OpsFileFinder.find_ops_files('test-cf-deployment')

        expect(ops).to include 'addons/ops.yml'
        expect(ops).to include 'addons/ops2.yml'
      end
    end
  end
end
