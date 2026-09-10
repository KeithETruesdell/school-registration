import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import {
  getStudents,
  getStudent,
  createStudent,
  updateStudent,
  deleteStudent,
  type Student,
} from "./students";

describe("students API", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  const mockStudent: Student = {
    id: 1,
    firstName: "John",
    lastName: "Doe",
    nickname: "JD",
    parentName: "Jane Doe",
    address1: "123 Main St",
    address2: "Apt 4",
    city: "Springfield",
    state: "IL",
    postalCode: "62701",
    email: "john@example.com",
    grade: "10",
  };

  describe("getStudents", () => {
    it("returns students on success", async () => {
      const students = [mockStudent];
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(students),
      } as Response);

      const result = await getStudents();

      expect(fetch).toHaveBeenCalledWith("/api/students");
      expect(result).toEqual(students);
    });

    it("throws on failure", async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
      } as Response);

      await expect(getStudents()).rejects.toThrow(
        "Unable to retrieve students",
      );
    });
  });

  describe("getStudent", () => {
    it("returns a student on success", async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockStudent),
      } as Response);

      const result = await getStudent(1);

      expect(fetch).toHaveBeenCalledWith("/api/students/1");
      expect(result).toEqual(mockStudent);
    });

    it("throws on failure", async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
      } as Response);

      await expect(getStudent(1)).rejects.toThrow(
        "Unable to retrieve student",
      );
    });
  });

  describe("createStudent", () => {
    it("creates and returns a student on success", async () => {
      const newStudent = {
        firstName: "John",
        lastName: "Doe",
        nickname: "JD",
        parentName: "Jane Doe",
        address1: "123 Main St",
        address2: "Apt 4",
        city: "Springfield",
        state: "IL",
        postalCode: "62701",
        email: "john@example.com",
        grade: "10",
      };

      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockStudent),
      } as Response);

      const result = await createStudent(newStudent);

      expect(fetch).toHaveBeenCalledWith("/api/students", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newStudent),
      });
      expect(result).toEqual(mockStudent);
    });

    it("throws on failure", async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
      } as Response);

      await expect(createStudent({} as Omit<Student, "id">)).rejects.toThrow(
        "Unable to create student",
      );
    });
  });

  describe("updateStudent", () => {
    it("updates and returns a student on success", async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockStudent),
      } as Response);

      const result = await updateStudent(1, mockStudent);

      expect(fetch).toHaveBeenCalledWith("/api/students/1", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(mockStudent),
      });
      expect(result).toEqual(mockStudent);
    });

    it("throws on failure", async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
      } as Response);

      await expect(updateStudent(1, mockStudent)).rejects.toThrow(
        "Unable to update student",
      );
    });
  });

  describe("deleteStudent", () => {
    it("deletes a student on success", async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
      } as Response);

      await deleteStudent(1);

      expect(fetch).toHaveBeenCalledWith("/api/students/1", {
        method: "DELETE",
      });
    });

    it("throws on failure", async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
      } as Response);

      await expect(deleteStudent(1)).rejects.toThrow(
        "Unable to delete student",
      );
    });
  });
});
