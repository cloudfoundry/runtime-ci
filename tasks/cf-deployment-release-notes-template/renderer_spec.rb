require 'rspec'
require_relative './renderer.rb'

describe 'Renderer' do
  describe '#render' do
    subject(:renderer) { Renderer.new }

    let(:release_update_1) do
      update = double('ReleaseUpdate')
      allow(update).to receive(:old_version) { '1.1.0' }
      allow(update).to receive(:new_version) { '1.3.0' }
      allow(update).to receive(:old_url)     { 'https://github.com/cloudfoundry/capi-release/releases/tag/v1.1.0' }
      allow(update).to receive(:new_url)     { 'https://github.com/cloudfoundry/capi-release/releases/tag/v1.3.0' }
      update
    end

    let(:release_update_2) do
      update = double('ReleaseUpdate')
      allow(update).to receive(:old_version) { '1.2.0' }
      allow(update).to receive(:new_version) { '1.4.0' }
      allow(update).to receive(:old_url)     { 'https://github.com/cloudfoundry/capi-release/releases/tag/v1.2.0' }
      allow(update).to receive(:new_url)     { 'https://github.com/cloudfoundry/capi-release/releases/tag/v1.4.0' }
      update
    end

    let(:release_updates) do
      updates = double('ReleaseUpdates')
      allow(updates).to receive(:each).and_yield('release-1', release_update_1).and_yield('release-2', release_update_2)
      updates
    end

    it 'includes a section header for Notices' do
      rendered_output = renderer.render(release_updates: release_updates)
      expect(rendered_output).to include ("## Notices\n\n")
    end

    it 'includes a section header for Manifest Updates' do
      rendered_output = renderer.render(release_updates: release_updates)
      expect(rendered_output).to include ("## Manifest Updates\n\n")
    end

    it 'includes a section header for Ops-files, as well as sub-headers for New and Updated Ops-files' do
      rendered_output = renderer.render(release_updates: release_updates)
      expect(rendered_output).to include (
<<-OPSFILES
## Ops-files
### New Ops-files
### Updated Ops-files

OPSFILES
      )
    end

    it 'includes a section header for Other Updates' do
      rendered_output = renderer.render(release_updates: release_updates)
      expect(rendered_output).to include ("## Other Updates\n\n")
    end

    it 'includes a section header for Release and Stemcell Updates' do
      rendered_output = renderer.render(release_updates: release_updates)
      expect(rendered_output).to include ("## Release Updates\n")
    end

    describe 'Release and stemcell table' do
      it 'includes a header' do
        expect(renderer.render(release_updates: release_updates)).to include(
<<-HEADER
| Release | Old Version | New Version | Release Notes |
| ------- | ----------- | ----------- | ------------- |
HEADER
        )
      end

      it 'includes an italicized disclaimer' do
        rendered_output = renderer.render(release_updates: release_updates)
        expect(rendered_output).to include ("_Warning: The Release Notes column only highlights noteworthy updates for each release bump. However, it is not exhaustive and we recommend you visit the actual release notes below before every upgrade._")
      end

      it 'places the disclaimer immediately after the section header' do
        rendered_output = renderer.render(release_updates: release_updates)
        expect(rendered_output).to include ("## Release Updates\n_Warning:")
      end

      it 'places the table header immediately after the disclaimer' do
        rendered_output = renderer.render(release_updates: release_updates)
        expect(rendered_output).to match (/Warning:.*\n| Release |/)
      end

      it 'shows the release name, old version, new version, and empty release notes for each release' do
        rendered_output = renderer.render(release_updates: release_updates)
        expect(rendered_output).to include('| release-1 | [1.1.0](https://github.com/cloudfoundry/capi-release/releases/tag/v1.1.0) | [1.3.0](https://github.com/cloudfoundry/capi-release/releases/tag/v1.3.0) | |')
        expect(rendered_output).to include('| release-2 | [1.2.0](https://github.com/cloudfoundry/capi-release/releases/tag/v1.2.0) | [1.4.0](https://github.com/cloudfoundry/capi-release/releases/tag/v1.4.0) | |')
      end

      context 'when some versions are nil' do
        let(:release_update_1) do
          update = double('ReleaseUpdate')
          allow(update).to receive(:old_version) { '1.1.0' }
          allow(update).to receive(:new_version) { nil }
          allow(update).to receive(:old_url)     { 'https://github.com/cloudfoundry/capi-release/releases/tag/v1.1.0' }
          update
        end

        let(:release_update_2) do
          update = double('ReleaseUpdate')
          allow(update).to receive(:old_version) { nil }
          allow(update).to receive(:new_version) { '1.4.0' }
          allow(update).to receive(:new_url)     { 'https://github.com/cloudfoundry/capi-release/releases/tag/v1.4.0' }
          update
        end

        it 'renders them as empty strings' do
          rendered_output = renderer.render(release_updates: release_updates)
          expect(rendered_output).to include('| release-1 | [1.1.0](https://github.com/cloudfoundry/capi-release/releases/tag/v1.1.0) |  |')
          expect(rendered_output).to include('| release-2 |  | [1.4.0](https://github.com/cloudfoundry/capi-release/releases/tag/v1.4.0) |')
        end
      end
    end
  end
end
