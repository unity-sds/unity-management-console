var text = `[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] HEAD points to '0bbe544506e8b58f2c6d2626ce901a724772fbeb'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] using github ref: refs/heads/main
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Loading revision from git directory
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] HEAD points to '0bbe544506e8b58f2c6d2626ce901a724772fbeb'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] using github ref: refs/heads/main
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Loading revision from git directory
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] HEAD points to '0bbe544506e8b58f2c6d2626ce901a724772fbeb'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] using github ref: refs/heads/main
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Loading revision from git directory
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] HEAD points to '0bbe544506e8b58f2c6d2626ce901a724772fbeb'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] using github ref: refs/heads/main
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] evaluating expression 'format('echo "The time was {0}"', steps.main_action_run.outputs.eksmeta)'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] expression 'format('echo "The time was {0}"', steps.main_action_run.outputs.eksmeta)' evaluated to '%!t(string=echo "The time was ")'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Wrote command 

echo "The time was "

 to 'workflow/7'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Writing entry to tarball workflow/7 len:22
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Extracting content to '/var/run/act'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute]   🐳  docker exec cmd=[bash --noprofile --norc -e -o pipefail /var/run/act/workflow/7] user= workdir=
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Exec command '[bash --noprofile --norc -e -o pipefail /var/run/act/workflow/7]'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Working directory '/home/ubuntu/unity-cs-infra'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute]   | The time was 
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute]   ✅  Success - Main Get the output time
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] skipping post step for 'Checkout': no action model available
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Removed container: 2150f0f7338ecc72a525b635e544d96eb340ce253dc470d49682300f0fa1359a
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG]   🐳  docker volume rm act-Development-Action-Not-For-Long-Term-Use-Only-Use-If-You-Kn-8ea346f3350bff32383107795ea0c298483ea47d23ac57f212f6539e752bfd24
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG]   🐳  docker volume rm act-Development-Action-Not-For-Long-Term-Use-Only-Use-If-You-Kn-8ea346f3350bff32383107795ea0c298483ea47d23ac57f212f6539e752bfd24-env
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] 🏁  Job succeeded
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Loading revision from git directory
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] HEAD points to '0bbe544506e8b58f2c6d2626ce901a724772fbeb'
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] using github ref: refs/heads/main
[Development Action Not For Long Term Use Only Use If You Know What You Are Doing/Parse Metadata and Execute] [DEBUG] Found revision: 0bbe544506e8b58f2c6d2626ce901a724772fbeb`
export const fetchline = (linenumber: number): string => {
    const lines = text.split("\n");
    if (linenumber >= 1 && linenumber <= lines.length) {
        return lines[linenumber - 1].trim() +'\n';
    }
    return "";
} 
