require 'rspec'
require_relative './release_changes.rb'

describe 'ReleaseUpdates' do
  describe '#load_change' do
    subject(:updates) { ReleaseUpdates.new }

    let(:version) { rand(10).to_s }
    let(:name) { 'capi-release' }
    let(:change) do
      [
        op,
        '[20]',
        {
          'name' => name,
          'version' => version
        }
      ]
    end

    context 'when the operation is "+"' do
      let(:op) { '+' }
      it 'saves the version as new_version' do
        updates.load_change(change)
        expect(updates.get_update_by_name(name).new_version).to eq(version)
      end
    end

    context 'when the operation is "-"' do
      let (:op) { '-' }
      it 'saves the version as old_version' do
        updates.load_change(change)
        expect(updates.get_update_by_name(name).old_version).to eq version
      end
    end

    context 'when a second change for the same release occurs' do
      let(:change1) do
        [
          '-',
          '[20]',
          {
            'name' => name,
            'version' => '26'
          }
        ]
      end

      let(:change2) do
        [
          '+',
          '[20]',
          {
            'name' => name,
            'version' => '27'
          }
        ]
      end

      it 'saves the old and new versions together' do
        subject.load_change(change1)
        subject.load_change(change2)
        expect(subject.get_update_by_name(name).new_version).to eq '27'
        expect(subject.get_update_by_name(name).old_version).to eq '26'
      end
    end
  end
end
