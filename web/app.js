const API_URL = "";


// =========================
// SAYFA GECISLERI
// =========================

function showRegister() {

    document.getElementById("loginPage").classList.add("hidden");

    document.getElementById("registerPage").classList.remove("hidden");

}


function showLogin() {

    document.getElementById("registerPage").classList.add("hidden");

    document.getElementById("loginPage").classList.remove("hidden");

}


// =========================
// KAYIT OL
// =========================

async function register() {

    const name =
        document.getElementById("registerName").value.trim();

    const email =
        document.getElementById("registerEmail").value.trim();

    const password =
        document.getElementById("registerPassword").value;

    const message =
        document.getElementById("registerMessage");


    if (!name || !email || !password) {

        message.innerText =
            "Tum alanlari doldur.";

        return;
    }


    try {

        const response = await fetch(
            `${API_URL}/users`,
            {
                method: "POST",

                headers: {
                    "Content-Type": "application/json"
                },

                body: JSON.stringify({
                    name: name,
                    email: email,
                    password: password
                })
            }
        );


        const data = await response.json();


        if (!response.ok) {

            message.innerText =
                data.error || "Kayit basarisiz.";

            return;
        }


        message.innerText =
            "Kayit basarili. Giris yapabilirsin.";


        document.getElementById("registerName").value = "";
        document.getElementById("registerEmail").value = "";
        document.getElementById("registerPassword").value = "";


        setTimeout(() => {

            showLogin();

        }, 1000);


    } catch (error) {

        console.error(error);

        message.innerText =
            "Sunucuya baglanilamadi.";

    }

}


// =========================
// GIRIS YAP
// =========================

async function login() {

    const email =
        document.getElementById("loginEmail").value.trim();

    const password =
        document.getElementById("loginPassword").value;

    const message =
        document.getElementById("loginMessage");


    if (!email || !password) {

        message.innerText =
            "E-posta ve sifre gir.";

        return;
    }


    try {

        const response = await fetch(
            `${API_URL}/login`,
            {
                method: "POST",

                headers: {
                    "Content-Type": "application/json"
                },

                body: JSON.stringify({
                    email: email,
                    password: password
                })
            }
        );


        const data = await response.json();


        if (!response.ok) {

            message.innerText =
                data.error || "Giris basarisiz.";

            return;
        }


        // JWT tokeni tarayicida sakla
        localStorage.setItem(
            "token",
            data.token
        );


        message.innerText =
            "Giris basarili.";


        showDashboard();


    } catch (error) {

        console.error(error);

        message.innerText =
            "Sunucuya baglanilamadi.";

    }

}


// =========================
// DASHBOARD
// =========================

function showDashboard() {

    document
        .getElementById("loginPage")
        .classList.add("hidden");


    document
        .getElementById("registerPage")
        .classList.add("hidden");


    document
        .getElementById("dashboardPage")
        .classList.remove("hidden");


    getFiles();

}


// =========================
// CIKIS
// =========================

function logout() {

    localStorage.removeItem("token");


    document
        .getElementById("dashboardPage")
        .classList.add("hidden");


    document
        .getElementById("loginPage")
        .classList.remove("hidden");


    document.getElementById("loginEmail").value = "";

    document.getElementById("loginPassword").value = "";

}


// =========================
// DOSYA YUKLE
// =========================

async function uploadFile() {

    const fileInput =
        document.getElementById("fileInput");

    const message =
        document.getElementById("uploadMessage");


    if (!fileInput.files.length) {

        message.innerText =
            "Lutfen bir dosya sec.";

        return;
    }


    const token =
        localStorage.getItem("token");


    if (!token) {

        message.innerText =
            "Once giris yapmalisin.";

        return;
    }


    const formData =
        new FormData();


    formData.append(
        "file",
        fileInput.files[0]
    );


    try {

        const response = await fetch(
            `${API_URL}/upload`,
            {
                method: "POST",

                headers: {
                    "Authorization":
                        `Bearer ${token}`
                },

                body: formData
            }
        );


        const data =
            await response.json();


        if (!response.ok) {

            message.innerText =
                data.error || "Dosya yuklenemedi.";

            return;
        }


        message.innerText =
            "Dosya basariyla yuklendi.";


        fileInput.value = "";


        getFiles();


    } catch (error) {

        console.error(error);

        message.innerText =
            "Sunucuya baglanilamadi.";

    }

}


// =========================
// DOSYALARIM
// =========================

async function getFiles() {

    const token =
        localStorage.getItem("token");


    if (!token) {
        return;
    }


    try {

        const response = await fetch(
            `${API_URL}/files`,
            {
                method: "GET",

                headers: {
                    "Authorization":
                        `Bearer ${token}`
                }
            }
        );


        const files =
            await response.json();


        if (!response.ok) {

            console.log(files);

            return;
        }


        const fileList =
            document.getElementById("fileList");


        if (!files || files.length === 0) {

            fileList.innerHTML =
                "Henuz dosya yuklenmedi.";

            return;
        }


        fileList.innerHTML = "";


        files.forEach(file => {

            const div =
                document.createElement("div");


            div.className =
                "fileItem";


            const downloadLink =
                `${window.location.origin}/download/${file.share_token}`;


            div.innerHTML = `

                <strong>
                    ${file.file_name}
                </strong>

                Boyut:
                ${file.file_size} byte

                <br>

                Son kullanma:
                ${new Date(file.expires_at)
                    .toLocaleString("tr-TR")}

                <div class="fileActions">

                    <button
                        onclick="window.open('${downloadLink}', '_blank')"
                    >
                        Indir
                    </button>


                    <button
                        onclick="copyLink('${downloadLink}')"
                    >
                        Linki Kopyala
                    </button>


                    <button
                        onclick="deleteFile(${file.id})"
                    >
                        Sil
                    </button>

                </div>
            `;


            fileList.appendChild(div);

        });


    } catch (error) {

        console.error(error);

    }

}


// =========================
// DOSYA SIL
// =========================

async function deleteFile(id) {

    const token =
        localStorage.getItem("token");


    if (!token) {

        alert("Once giris yapmalisin.");

        return;
    }


    const confirmDelete =
        confirm(
            "Bu dosyayi silmek istedigine emin misin?"
        );


    if (!confirmDelete) {
        return;
    }


    try {

        const response = await fetch(
            `${API_URL}/files/${id}`,
            {
                method: "DELETE",

                headers: {
                    "Authorization":
                        `Bearer ${token}`
                }
            }
        );


        const data =
            await response.json();


        if (!response.ok) {

            alert(
                data.error ||
                "Dosya silinemedi."
            );

            return;
        }


        alert(
            data.message ||
            "Dosya basariyla silindi."
        );


        getFiles();


    } catch (error) {

        console.error(error);

        alert(
            "Sunucuya baglanilamadi."
        );

    }

}


// =========================
// LINK KOPYALA
// =========================

async function copyLink(link) {

    try {

        await navigator.clipboard.writeText(link);

        alert(
            "Paylasim linki kopyalandi."
        );

    } catch (error) {

        console.error(error);

        alert(
            "Link kopyalanamadi."
        );

    }

}


// =========================
// SAYFA ACILINCA
// =========================

window.onload = function () {

    const token =
        localStorage.getItem("token");


    if (token) {

        showDashboard();

    }

};