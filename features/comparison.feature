Feature:

  Scenario: Simple string to string comparison, equal
    Given the program:
    """
    when
      "x" == "x"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple string to string comparison, not equal
    Given the program:
    """
    when
      "x" == "y"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Simple variable to string comparison, equal
    Given the program:
    """
    when
      var(a) == "x"
    then
      score(x) = 1
    done
    """
    And variables:
      | Name | Value |
      | a    | x     |
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple variable to string comparison, not equal
    Given the program:
    """
    when
      var(a) == "y"
    then
      score(x) = 1
    done
    """
    And variables:
      | Name | Value |
      | a    | x     |
    When the program is run
    Then the score output is empty

  Scenario: Simple score to int comparison, equal
    Given the program:
    """
    when
      score(x) == 0
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple score to int comparison, not equal
    Given the program:
    """
    when
      score(x) == 1
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Simple int to int greater than
    Given the program:
    """
    when
      1 > 0
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple int to int not greater than
    Given the program:
    """
    when
      0 > 1
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Simple int to int greater than or equal
    Given the program:
    """
    when
      1 >= 1
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple int to int not greater than or equal
    Given the program:
    """
    when
      0 >= 1
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Simple int to int less than
    Given the program:
    """
    when
      0 < 1
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple int to int not less than
    Given the program:
    """
    when
      1 < 0
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Simple int to int less than or equal
    Given the program:
    """
    when
      1 <= 1
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple int to int not less than or equal
    Given the program:
    """
    when
      1 <= 0
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Simple int to int not equals
    Given the program:
    """
    when
      1 != 2
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple int to int not equals but equal
    Given the program:
    """
    when
      1 != 1
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Simple string to string not equals
    Given the program:
    """
    when
      "x" != "y"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple string to var not equals
    Given the program:
    """
    when
      "y" != var(a)
    then
      score(x) = 1
    done
    """
    And variables:
      | Name | Value |
      | a    | x     |
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple string to string not equals but equal
    Given the program:
    """
    when
      "x" != "x"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Simple string to string contains
    Given the program:
    """
    when
      "xyz" contains "x"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple string to string not equals but equal
    Given the program:
    """
    when
      "xyz" contains "a"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty