const sendMail = () => {
    let name = document.getElementById("name-input").value;
    let email = document.getElementById("email-input").value;
    let phone = document.getElementById("phone-input").value;
    let subject = document.getElementById("subject-input").value;
    let message = document.getElementById("message-input").value;

    if (name === "") {
        return alert("Name is required");
    }
    if (email === "") {
        return alert("email is required");
    }
    if (phone === "") {
        return alert("Phone is required");
    }
    if (subject === "") {
        return alert("Subject is required");
    }

    const emailReciever = "qidamqolby@gmail.com";

    const a = document.createElement("a");

    a.href = `mailto:${emailReciever}?subject=${subject}&body= Hello, my name is ${name}, ${message}, please contact me at ${phone}`;
    console.log(a);
    a.click();

    alert("pesan anda telah terkirim");
};
