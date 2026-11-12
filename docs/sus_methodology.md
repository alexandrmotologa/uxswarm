# System Usability Scale (SUS) methodology

uxswarm calculates the standard System Usability Scale (Brooke, 1996) score for each synthetic session and aggregates them into a swarm benchmark.

## Standard questionnaire

The evaluation uses the ten standard SUS statements rated on a 1 to 5 scale (1 = Strongly Disagree, 5 = Strongly Agree):

1. I think that I would like to use this system frequently.
2. I found the system unnecessarily complex.
3. I thought the system was easy to use.
4. I think that I would need the support of a technical person to be able to use this system.
5. I found the various functions in this system were well integrated.
6. I thought there was too much inconsistency in this system.
7. I would imagine that most people would learn to use this system very quickly.
8. I found the system very cumbersome to use.
9. I felt very confident using the system.
10. I needed to learn a lot of things before I could get going with this system.

## Score calculation

Let R[i] represent the rating (1 to 5) for item i:

- For odd items (1, 3, 5, 7, 9): S[i] = R[i] - 1
- For even items (2, 4, 6, 8, 10): S[i] = 5 - R[i]

Total SUS score = sum(S[1..10]) * 2.5

The resulting score ranges from 0 to 100.

## Grading scale

- Score >= 80.3: Grade A (High usability, low user friction)
- 68.0 <= Score < 80.3: Grade B or C (Acceptable usability, minor friction points)
- 51.0 <= Score < 68.0: Grade D (Substandard usability, significant friction)
- Score < 51.0: Grade F (Severe usability barriers, high abandonment risk)

## Friction Index (FI)

In addition to post-session questionnaires, telemetry during live navigation calculates the Friction Index:

FI = (w1 * (T_hesitation / T_total)) + (w2 * N_misclicks) + (w3 * N_backtracks)

Default weights:
- w1 = 0.35 (Hesitation proportion)
- w2 = 0.40 (Dead click count)
- w3 = 0.25 (Backtrack loops)

An FI value greater than 2.5 indicates high friction and marks the tested flow for redesign.
