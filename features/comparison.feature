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

  Scenario: Simple string to string contains failure
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

  Scenario: Simple string to string does not contain
    Given the program:
    """
    when
      "xyz" does not contain "a"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple string regexp match
    Given the program:
    """
    when
        "xyz" matches /y/
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: String regexp match with dot escaping
    Given the program:
    """
    when
        "." matches /^\.$/
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: String regexp match with slash escaping
    Given the program:
    """
    when
        "/" matches /^\/$/
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: String regexp match with character class
    Given the program:
    """
    when
        "1" matches /^\d$/
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Simple string regexp match failure
    Given the program:
    """
    when
        "xyz" matches /a/
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Simple string regexp does not match
    Given the program:
    """
    when
        "xyz" does not match /a/
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: String in list
    Given the program:
    """
    when
        "y" in ["x", "y", "z"]
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Var in list
    Given the program:
    """
    when
        var(a) in ["x", "y", "z"]
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

  Scenario: String in list containing mixed values
    Given the program:
    """
    when
        "y" in ["x", var(a), "z"]
    then
        score(x) = 1
    done
    """
    And variables:
      | Name | Value |
      | a    | y     |
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: String in list failure
    Given the program:
    """
    when
        "a" in ["x", "y", "z"]
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: String not in list
    Given the program:
    """
    when
        "a" not in ["x", "y", "z"]
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Var not in list
    Given the program:
    """
    when
        var(a) not in ["x", "y", "z"]
    then
        score(x) = 1
    done
    """
    And variables:
      | Name | Value |
      | a    | a     |
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: String not in list containing mixed values
    Given the program:
    """
    when
        "y" not in ["x", var(a), "z"]
    then
        score(x) = 1
    done
    """
    And variables:
      | Name | Value |
      | a    | a     |
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: String not in list failure
    Given the program:
    """
    when
        "z" not in ["x", "y", "z"]
    then
        score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty
