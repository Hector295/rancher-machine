package provision

import (
	"fmt"
	"os"

	"github.com/rancher/machine/libmachine/provision/pkgaction"
)

func WithCustomScript(provisioner Provisioner, customScriptPath, hostname string) error {
	if provisioner == nil {
		return nil
	}

	if hostname == "" {
		hostname = provisioner.GetDriver().GetMachineName()
	}

	if err := provisioner.SetHostname(hostname); err != nil {
		return err
	}

	for _, pkg := range provisioner.GetPackages() {
		if err := provisioner.Package(pkg, pkgaction.Install); err != nil {
			return err
		}
	}

	customScriptContents, err := os.ReadFile(customScriptPath)
	if err != nil {
		return fmt.Errorf("unable to read file %s: %v", customScriptPath, err)
	}

	if output, err := provisioner.SSHCommand(fmt.Sprintf("cat <<'OEOF' >/tmp/install_script.sh\n%s\nOEOF", string(customScriptContents))); err != nil {
		return fmt.Errorf("error uploading custom script: output: %s, error: %s", output, err)
	}
	if output, err := provisioner.SSHCommand("sudo sh /tmp/install_script.sh"); err != nil {
		return fmt.Errorf("error running custom script: output: %s, error: %s", output, err)
	}

	return nil
}
