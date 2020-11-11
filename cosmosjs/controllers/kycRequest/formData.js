var FormData = require('form-data');
var form = new FormData();

module.exports = {
    generateKycData: (data) => {
        form.append("oracle_script_name", data.oracle_script_name)
        form.append("from", data.from)
        form.append("chain_id", data.chain_id)
        form.append("image", './.uploads' + data.img_name)
        form.append("input", data.input)
        form.append("expected_output", data.expected_output)
        form.append("fees", data.fees)
        form.append("validator_count", data.validator_count)
        return form
    }
}