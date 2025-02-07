# Easy Japanese
Frontend
* Navigation between different modules
* Smooth animations

Backend
* Basics
  * Log in and verification.
  * Store information about users.
* Functions
  * Vocabulary. Get a random word from vocabulary. Recently chosen words have a lower probability of being chosen again. The weight of a word depends on how well the user has mastered it. There is a lower bound and an upper bound for weight. Current thought is to maintain several sets of words, say master, familiar, known, unknown. Words in master will be chosen at a very low probability.  
  * Grammars (expression). Based on translation. LLM will be used to grade the answer of the user. Similar trick may be used.
  * Grammars (comprehension). Simply choose a random sentence from the database. 
* Advanced
  * User data are saved every several seconds.