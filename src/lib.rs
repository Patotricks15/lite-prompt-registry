//! Immutable prompt-version governance without network or database overhead.

use std::collections::HashMap;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Prompt {
    pub id: u64,
    pub name: String,
    pub description: String,
    pub versions: Vec<Version>,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Version {
    pub number: u32,
    pub template: String,
    pub author: String,
    pub status: Status,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Status { Draft, InReview, Approved, RolledOut }

#[derive(Default)]
pub struct Registry { prompts: HashMap<u64, Prompt>, next_id: u64 }

impl Registry {
    pub fn create_prompt(&mut self, name: impl Into<String>, description: impl Into<String>) -> u64 {
        self.next_id += 1;
        self.prompts.insert(self.next_id, Prompt { id: self.next_id, name: name.into(), description: description.into(), versions: Vec::new() });
        self.next_id
    }

    pub fn create_version(&mut self, id: u64, template: impl Into<String>, author: impl Into<String>) -> Result<u32, Error> {
        let prompt = self.prompts.get_mut(&id).ok_or(Error::PromptNotFound)?;
        let number = prompt.versions.len() as u32 + 1;
        prompt.versions.push(Version { number, template: template.into(), author: author.into(), status: Status::Draft });
        Ok(number)
    }

    pub fn get(&self, id: u64) -> Option<&Prompt> { self.prompts.get(&id) }
}

#[derive(Debug, thiserror::Error)]
pub enum Error { #[error("prompt not found")] PromptNotFound }

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn versions_are_sequential() {
        let mut registry = Registry::default();
        let id = registry.create_prompt("welcome", "");
        assert_eq!(registry.create_version(id, "Hello", "a").unwrap(), 1);
    }
}