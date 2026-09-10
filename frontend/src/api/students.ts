export interface Student {
  id: number;
  firstName: string;
  lastName: string;
  nickname: string;
  parentName: string;
  address1: string;
  address2: string;
  city: string;
  state: string;
  postalCode: string;
  email: string;
  grade: string;
}

export async function getStudents(): Promise<Student[]> {
  const response =
    await fetch("/api/students");

  if (!response.ok) {
    throw new Error(
      "Unable to retrieve students",
    );
  }

  return response.json();
}

export async function getStudent(
  id: number,
): Promise<Student> {
  const response =
    await fetch(`/api/students/${id}`);

  if (!response.ok) {
    throw new Error(
      "Unable to retrieve student",
    );
  }

  return response.json();
}
