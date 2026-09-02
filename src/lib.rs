//! Immutable prompt-version governance without network or database overhead.

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs;
use std::path::Path;

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Prompt {
    #[serde(default)]
    pub id: String,
    pub name: String,
    #[serde(default)]
    pub description: String,
    #[serde(default)]
    pub versions: Vec<Version>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
pub struct Version {
    pub version: u32,
    pub template: String,
    #[serde(default)]
    pub author: String,
    #[serde(default)]
    pub status: Status,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Default, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum Status {
    #[default]
    Draft,
    InReview,
    Approved,
    RolledOut,
}

#[derive(Debug, Deserialize)]
struct PromptYamlDoc {
    prompts: Vec<Prompt>,
}

#[derive(Default, Debug, Clone, Serialize, Deserialize)]
pub struct Registry {
    prompts: HashMap<String, Prompt>,
}

impl Registry {
    /// Loads and validates prompt definitions directly from a YAML or JSON file.
    pub fn from_file(path: impl AsRef<Path>) -> Result<Self, Box<dyn std::error::Error>> {
        let content = fs::read_to_string(path.as_ref())?;
        let doc: PromptYamlDoc = serde_yaml::from_str(&content)?;
        let mut map = HashMap::new();
        for p in doc.prompts {
            map.insert(p.id.clone(), p);
        }
        Ok(Self { prompts: map })
    }

    pub fn register(&mut self, id: impl Into<String>, name: impl Into<String>, description: impl Into<String>) {
        let id_str = id.into();
        self.prompts.insert(
            id_str.clone(),
            Prompt {
                id: id_str,
                name: name.into(),
                description: description.into(),
                versions: Vec::new(),
            },
        );
    }

    pub fn add_version(
        &mut self,
        id: &str,
        template: impl Into<String>,
        author: impl Into<String>,
        status: Status,
    ) -> Result<u32, Error> {
        let prompt = self.prompts.get_mut(id).ok_or(Error::PromptNotFound)?;
        let version_num = prompt.versions.len() as u32 + 1;
        prompt.versions.push(Version {
            version: version_num,
            template: template.into(),
            author: author.into(),
            status,
        });
        Ok(version_num)
    }

    pub fn get(&self, id: &str) -> Option<&Prompt> {
        self.prompts.get(id)
    }

    /// Renders a prompt version by substituting `{key}` variables.
    pub fn render(&self, id: &str, version: Option<u32>, variables: &HashMap<&str, &str>) -> Result<String, Error> {
        let prompt = self.prompts.get(id).ok_or(Error::PromptNotFound)?;
        let ver = match version {
            Some(v) => prompt.versions.iter().find(|x| x.version == v).ok_or(Error::VersionNotFound)?,
            None => prompt.versions.last().ok_or(Error::NoVersionsAvailable)?,
        };

        let mut output = ver.template.clone();
        for (k, v) in variables {
            output = output.replace(&format!("{{{}}}", k), v);
        }
        Ok(output)
    }
}

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("prompt not found")]
    PromptNotFound,
    #[error("version not found")]
    VersionNotFound,
    #[error("no versions available for prompt")]
    NoVersionsAvailable,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn loads_and_renders_yaml_prompts() {
        let registry = Registry::from_file("examples/prompts.yaml").unwrap();
        let prompt = registry.get("customer_support").expect("prompt should exist");
        assert_eq!(prompt.versions.len(), 2);

        let mut vars = HashMap::new();
        vars.insert("company", "CloudInc");
        vars.insert("user_query", "How to pay?");

        let rendered = registry.render("customer_support", Some(1), &vars).unwrap();
        assert!(rendered.contains("CloudInc"));
        assert!(rendered.contains("How to pay?"));
    }
}