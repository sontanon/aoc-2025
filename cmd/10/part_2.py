from pathlib import Path
from dataclasses import dataclass
from pulp import LpProblem, LpMinimize, LpVariable, lpSum, value


@dataclass
class Button:
    indices: list[int]

    def to_column(self, length: int) -> list[int]:
        column = [0] * length
        for index in self.indices:
            column[index] = 1
        return column


@dataclass
class ActionSpace:
    buttons: list[Button]
    goal: list[int]

    def to_matrix(self) -> list[list[int]]:
        matrix = [[0] * len(self.buttons) for _ in range(len(self.goal))]
        for j, button in enumerate(self.buttons):
            column = button.to_column(len(self.goal))
            for i in range(len(self.goal)):
                matrix[i][j] = column[i]

        return matrix

    def solve(self) -> list[int]:
        A = self.to_matrix()
        y = self.goal

        prob = LpProblem("Minimize_Sum_x", LpMinimize)
        x = [
            LpVariable(f"x_{i}", lowBound=0, cat="Integer")
            for i in range(len(self.buttons))
        ]
        prob += lpSum(x)

        for i in range(len(y)):
            prob += lpSum(A[i][j] * x[j] for j in range(len(self.buttons))) == y[i]

        prob.solve()
        return [int(value(var)) for var in x]


def parse_button(input: str) -> Button:
    if len(input) < 3:
        raise ValueError("input is too short to be valid")

    if input[0] != "(" or input[-1] != ")":
        raise ValueError("input is not enclosed in parentheses")

    sub_fields = input[1:-1].split(",")
    indices = [0 for _ in range(len(sub_fields))]
    for i, sub_field in enumerate(sub_fields):
        try:
            indices[i] = int(sub_field)
        except ValueError as e:
            raise ValueError(
                f"error processing field '{sub_field}' at index {i}"
            ) from e
    return Button(indices)


def parse_joltage(input: str) -> list[int]:
    if len(input) < 3:
        raise ValueError("input is too short to be valid")
    if input[0] != "{" or input[-1] != "}":
        raise ValueError("input is not enclosed in curly braces")
    sub_fields = input[1:-1].split(",")
    joltage = [0 for _ in range(len(sub_fields))]
    for i, sub_field in enumerate(sub_fields):
        try:
            joltage[i] = int(sub_field)
        except ValueError as e:
            raise ValueError(
                f"error processing field '{sub_field}' at index {i}"
            ) from e
    return joltage


def parse_action_space(line: str) -> ActionSpace:
    fields = line.strip().split()

    buttons = [Button([]) for _ in range(len(fields) - 2)]
    for i, field in enumerate(fields[1:-1]):
        buttons[i] = parse_button(field)

    goal = parse_joltage(fields[-1])
    return ActionSpace(buttons, goal)


def main():
    input_path = Path(__file__).parent / "input.txt"
    input_lines = input_path.read_text().strip().splitlines()

    total = 0
    for line in input_lines:
        action_space = parse_action_space(line)
        try:
            result = action_space.solve()
            total += sum(result)
        except Exception as e:
            print(f"Error solving action space for line: {line}")
            print(e)

    print(f"Total sum of button presses: {total}")


if __name__ == "__main__":
    main()
