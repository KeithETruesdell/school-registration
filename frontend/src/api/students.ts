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

export async function createStudent(
  student: Omit<Student, "id">,
): Promise<Student> {
  const response = await fetch(
    "/api/students",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(student),
    },
  );

  if (!response.ok) {
    throw new Error(
      "Unable to create student",
    );
  }

  return response.json();
}

export async function updateStudent(
  id: number,
  student: Student,
): Promise<Student> {
  const response = await fetch(
    `/api/students/${id}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(student),
    },
  );

  if (!response.ok) {
    throw new Error(
      "Unable to update student",
    );
  }

  return response.json();
}

export async function deleteStudent(
  id: number,
): Promise<void> {
  const response = await fetch(
    `/api/students/${id}`,
    {
      method: "DELETE",
    },
  );

  if (!response.ok) {
    throw new Error(
      "Unable to delete student",
    );
  }
}
