## Business Requirements

The enterprise wants to develop features for their existing language portal and has shown interest in using AI to powered said features. The reason for this is that the students need various tools to independently learn and practice, giving teachers the time they need to properly assest every student. The goal with these feature is to enable students to become autodidacts with the help of AI, as for the teachers, the goal is to give them the space to properly assest each student and prepare classes around the needs of the students. 

The enterprise wants to keep cost an minimum while protecting the private data of each student, the LLMs that will be used should provide a good balance between costs, security and privacy of the students information.

## Functional Requirements

- The LLM should give suggestions for improvements on grammar and context based on the input provider by the student.
- The LLM should provide a break down, explaining grammar structure, of the sentence provided by the student.
- The LLM should provide examples in different levels of complexity  of the sentence provided by the student.

## Non Functional Requirements

- The use of a RAG is recommended, this will directly impact the precision of the answer of the LLM.
- The use of a cache for the user queries will help improve the time of the RAG itself.

## Assumptions

- We're assuming that an open-source on premise LLM will provide a lower cost than a cloud base LLM.

## Risks

- Cybersecurity, not properly following guidelines and best practices can result in a breach and lost of the system or data.
- Prompt injection, students can find ways to exploit the LLMs for other things outside the tools and features offered to them.

## Constraint

- To keep a low cost there are some limits on the hardware and technologies available for the project.

## Data Strategy

- For the writing practicing app feature students inputs will be run through a RAG solution will be put in place to enhance the precision, this solution requires a vector type database where an index will be persisted. No personal information of the students will be collected in this dataset to keep risk of data leakage low.

## Model Selection and Development

The solution will be self hosted, open source, the feature in question is a text-to-text feature. The expectancy is that a 5 billion parameter model should suffice and some fine tunning will be required to get better response specially when inquiring complex sentences.

