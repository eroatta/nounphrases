# phrase-finder

## Introduction

_phrase-finder_ is an implementation for the  "probabilistic model for automatically extracting English noun phrases without part-of-speech tagging or any syntactic analysis", proposed by Feng and Croft on their paper called "Probabilistic techniques for phrase extraction".
As they stated, the technique is based on a **Markov model**, whose initial parameters are estimated by a phrase lookup program with a phrase dictionary, then optimized by a set of **maximum entropy** (ME) parameters for a set of morphological features.
Using the **Viterbi algorithm** with the trained Markov model, the program can dynamically extract noun phrases from input text.

## Usage

WIP

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
