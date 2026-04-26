use ed25519_dalek::{Signature, Signer, SigningKey, Verifier, VerifyingKey};
use prost::Message;

#[derive(Clone, PartialEq, Message)]
pub struct UAPEnvelope {
    #[prost(string, tag = "1")]
    pub version: String,
    #[prost(string, tag = "2")]
    pub from_agent: String,
    #[prost(string, tag = "3")]
    pub to_agent: String,
    #[prost(enumeration = "MessageType", tag = "4")]
    pub r#type: i32,
    #[prost(bytes = "vec", tag = "5")]
    pub payload: Vec<u8>,
    #[prost(bytes = "vec", tag = "6")]
    pub signature: Vec<u8>,
    #[prost(int64, tag = "7")]
    pub timestamp: i64,
    #[prost(string, tag = "8")]
    pub nonce: String,
}

#[derive(Clone, Copy, Debug, PartialEq, Eq, prost::Enumeration)]
#[repr(i32)]
pub enum MessageType {
    CapabilityAnnounce = 0,
    TaskRequest = 1,
    TaskAccept = 2,
    TaskResult = 3,
    KnowledgeShare = 4,
    HelpRequest = 5,
    EvolutionProposal = 6,
    ConsensusVote = 7,
}

impl UAPEnvelope {
    pub fn signed(
        mut envelope: UAPEnvelope,
        signer: &SigningKey,
    ) -> anyhow::Result<UAPEnvelope> {
        envelope.signature.clear();
        let bytes = envelope.encode_to_vec();
        let signature = signer.sign(&bytes);
        envelope.signature = signature.to_bytes().to_vec();
        Ok(envelope)
    }

    pub fn verify(&self, public_key: &VerifyingKey) -> anyhow::Result<()> {
        let mut cloned = self.clone();
        let sig = Signature::from_slice(&cloned.signature)?;
        cloned.signature.clear();
        let bytes = cloned.encode_to_vec();
        public_key.verify(&bytes, &sig)?;
        Ok(())
    }
}
