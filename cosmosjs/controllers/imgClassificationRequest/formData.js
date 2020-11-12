const FormData = require('form-data');
const form = new FormData();
const fs = require('fs');

module.exports = {
    generateClassificationData: (data) => {
        form.append("oracle_script_name", data.oracle_script_name)
        form.append("from", data.from)
        form.append("chain_id", data.chain_id)
        form.append("image", fs.createReadStream(data.img_path))
        form.append("input", data.input)
        form.append("expected_output", data.expected_output)
        form.append("fees", data.fees)
        form.append("validator_count", data.validator_count)
        return form
    }
}