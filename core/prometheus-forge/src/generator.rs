use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ToolSpec {
    pub name: String,
    pub description: String,
    pub inputs: Vec<String>,
    pub outputs: Vec<String>,
}

pub struct ToolGenerator;

impl ToolGenerator {
    pub fn generate_stub(&self, spec: &ToolSpec) -> String {
        format!(
            "def {}({}):\n    \"\"\"{}\"\"\"\n    raise NotImplementedError\n",
            spec.name,
            spec.inputs.join(", "),
            spec.description
        )
    }
}
