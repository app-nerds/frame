import { fetcher } from "../http/fetcher.js";

export class MemberService {
  spinnerEl;

  constructor(spinnerEl) {
    this.spinnerEl = spinnerEl;
  }

  async getCurrentMember() {
    let options = {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    };

    let response = await fetcher(`/api/member/current`, options, this.spinnerEl);
    let result = await response.json();
    return result;
  }
}
