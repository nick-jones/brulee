Feature:

  Scenario: And, all conditions resolve true
    Given the program:
    """
    when
      "x" == "x" and "y" == "y"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: And, second condition fails
    Given the program:
    """
    when
      "x" == "x" and "y" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: And, first condition fails
    Given the program:
    """
    when
      "x" == "y" and "z" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: And, 3 conditions
    Given the program:
    """
    when
      "x" == "x" and "y" == "y" and "z" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: And, 3 conditions, middle fails
    Given the program:
    """
    when
      "x" == "x" and "y" == "a" and "z" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Or, all conditions resolve true
    Given the program:
    """
    when
      "x" == "x" or "y" == "y"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Or, first condition resolve true
    Given the program:
    """
    when
      "x" == "x" or "y" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Or, second condition resolve true
    Given the program:
    """
    when
      "x" == "z" or "y" == "y"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Or, 3 conditions
    Given the program:
    """
    when
      "x" == "z" or "y" == "y" or "z" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Mixed
    Given the program:
    """
    when
      "x" == "z" or "y" == "y" and "z" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Mixed, first and second condition fails
    Given the program:
    """
    when
      "x" == "z" or "y" == "z" and "z" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Mixed, second condition fails
    Given the program:
    """
    when
      "x" == "x" or "y" == "z" and "z" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Mixed, first condition fails
    Given the program:
    """
    when
      "x" == "y" or "y" == "y" and "z" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Grouped conditions
    Given the program:
    """
    when
      ("x" == "y" or "y" == "y") and "z" == "z"
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Grouped conditions (2)
    Given the program:
    """
    when
      "z" == "z" and ("y" == "y" or "x" == "y")
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Grouped conditions (3)
    Given the program:
    """
    when
      "z" == "z" and ("y" == "z" or "x" == "y")
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Grouped conditions (4)
    Given the program:
    """
    when
      "y" == "x" and ("y" == "y" or "z" == "z")
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: Grouped conditions (5)
    Given the program:
    """
    when
      "y" == "x" or ("y" == "y" and "z" == "z")
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Grouped conditions (6)
    Given the program:
    """
    when
      "y" == "y" or ("y" == "y" and "z" == "x")
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: Grouped conditions (7)
    Given the program:
    """
    when
      "x" == "y" or ("y" == "y" and "z" == "x")
    then
      score(x) = 1
    done
    """
    When the program is run
    Then the score output is empty