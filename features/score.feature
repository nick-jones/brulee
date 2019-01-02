Feature:

  Scenario: Score adjustment
    Given the program:
    """
    score(x) = 1
    score(x) += 1
    score(y) = 10
    score(y) -= 1
    score(z) = score(y)
    score(z) += 1
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 2     |
      | y    | 9     |
      | z    | 10    |