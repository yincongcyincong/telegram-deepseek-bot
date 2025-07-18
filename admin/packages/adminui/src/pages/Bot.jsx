import React, { useEffect, useState } from "react";
import Modal from "../components/Modal";
import Pagination from "../components/Pagination";
import ConfigForm from "./ConfigForm";
import Toast from "../components/Toast";
import ConfirmModal from "../components/ConfirmModal";

function Bots() {
    const [bots, setBots] = useState([]);
    const [search, setSearch] = useState("");
    const [modalVisible, setModalVisible] = useState(false);
    const [editingBot, setEditingBot] = useState(null);
    const [form, setForm] = useState({
        id: 0,
        address: "",
        crt_file: "",
        key_file: "",
        ca_file: "",
    });

    const [rawConfigVisible, setRawConfigVisible] = useState(false);
    const [structuredConfigVisible, setStructuredConfigVisible] = useState(false);
    const [rawConfigText, setRawConfigText] = useState("");
    const [selectId, setSelectId] = useState(0);

    const [page, setPage] = useState(1);
    const [pageSize] = useState(10);
    const [total, setTotal] = useState(0);

    const [toast, setToast] = useState({ show: false, message: "", type: "error" });
    const [confirmVisible, setConfirmVisible] = useState(false);
    const [botToDelete, setBotToDelete] = useState(null);

    const showToast = (message, type = "error") => {
        setToast({ show: true, message, type });
    };

    useEffect(() => {
        fetchBots();
    }, [page]);

    const fetchBots = async () => {
        try {
            const res = await fetch(
                `/bot/list?page=${page}&page_size=${pageSize}&address=${encodeURIComponent(search)}`
            );
            const data = await res.json();
            if (data.code !== 0) {
                showToast(data.message || "Failed to fetch bots");
                return;
            }
            setBots(data.data.list);
            setTotal(data.data.total);
        } catch (err) {
            showToast("Request error: " + err.message);
        }
    };

    const handleSearch = () => {
        setPage(1);
        fetchBots();
    };

    const handleAddClick = () => {
        setForm({ id: 0, address: "", crt_file: "", key_file: "", ca_file: "" });
        setEditingBot(null);
        setModalVisible(true);
    };

    const handleEditClick = (bot) => {
        setForm({
            id: bot.id,
            address: bot.address,
            crt_file: bot.crt_file,
            key_file: bot.key_file,
            ca_file: bot.ca_file,
        });
        setEditingBot(bot);
        setModalVisible(true);
    };

    const handleDeleteClick = (botId) => {
        setBotToDelete(botId);
        setConfirmVisible(true);
    };

    const cancelDelete = () => {
        setBotToDelete(null);
        setConfirmVisible(false);
    };

    const confirmDelete = async () => {
        if (!botToDelete) return;
        try {
            const res = await fetch(`/bot/delete?id=${botToDelete}`, { method: "DELETE" });
            const data = await res.json();
            if (data.code !== 0) {
                showToast(data.message || "Failed to delete bot");
                return;
            }
            showToast("Bot deleted", "success");
            setConfirmVisible(false);
            setBotToDelete(null);
            await fetchBots();
        } catch (error) {
            showToast("Request error: " + error.message);
        }
    };

    const handleSave = async () => {
        try {
            const url = editingBot ? "/bot/update" : "/bot/create";
            const res = await fetch(url, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(form),
            });
            const data = await res.json();
            if (data.code !== 0) {
                showToast(data.message || "Failed to save bot");
                return;
            }
            await fetchBots();
            setModalVisible(false);
        } catch (err) {
            showToast("Request error: " + err.message);
        }
    };

    const handleShowRawConfig = async (botId) => {
        try {
            const res = await fetch(`/bot/command/get?id=${botId}`);
            const data = await res.json();
            if (data.code !== 0) {
                showToast(data.message || "Failed to fetch command config");
                return;
            }
            setRawConfigText(data.data);
            setRawConfigVisible(true);
        } catch (err) {
            showToast("Request error: " + err.message);
        }
    };

    const handleShowStructuredConfig = (botId) => {
        setStructuredConfigVisible(true);
        setSelectId(botId);
    };

    const handlePageChange = (newPage) => {
        setPage(newPage);
    };

    return (
        <div className="p-6 bg-gray-100 min-h-screen">
            {toast.show && (
                <Toast
                    message={toast.message}
                    type={toast.type}
                    onClose={() => setToast({ ...toast, show: false })}
                />
            )}

            <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-bold text-gray-800">Bot Management</h2>
                <button
                    onClick={handleAddClick}
                    className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded shadow"
                >
                    + Add Bot
                </button>
            </div>

            <div className="flex mb-4 space-x-2">
                <input
                    type="text"
                    placeholder="Search by address"
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    className="w-full sm:w-64 px-4 py-2 border border-gray-300 rounded shadow-sm focus:outline-none focus:ring focus:border-blue-400"
                />
                <button
                    onClick={handleSearch}
                    className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                >
                    Search
                </button>
            </div>

            <div className="overflow-x-auto rounded-lg shadow">
                <table className="min-w-full bg-white divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                    <tr>
                        {[
                            "ID",
                            "Address",
                            "Status",
                            "Create Time",
                            "Update Time",
                            "Actions",
                        ].map((title) => (
                            <th
                                key={title}
                                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                            >
                                {title}
                            </th>
                        ))}
                    </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                    {bots.map((bot) => (
                        <tr key={bot.id} className="hover:bg-gray-50">
                            <td className="px-6 py-4 text-sm text-gray-800">{bot.id}</td>
                            <td className="px-6 py-4 text-sm text-gray-800">{bot.address}</td>
                            <td className="px-6 py-4 text-sm text-gray-800">{bot.status}</td>
                            <td className="px-6 py-4 text-sm text-gray-600">
                                {new Date(bot.create_time * 1000).toLocaleString()}
                            </td>
                            <td className="px-6 py-4 text-sm text-gray-600">
                                {new Date(bot.update_time * 1000).toLocaleString()}
                            </td>
                            <td className="px-6 py-4 space-x-2 text-sm">
                                <button
                                    onClick={() => handleEditClick(bot)}
                                    className="text-blue-600 hover:underline"
                                >
                                    Edit
                                </button>
                                <button
                                    onClick={() => handleShowRawConfig(bot.id)}
                                    className="text-purple-600 hover:underline"
                                >
                                    Command
                                </button>
                                <button
                                    onClick={() => handleShowStructuredConfig(bot.id)}
                                    className="text-green-600 hover:underline"
                                >
                                    Config
                                </button>
                                <button
                                    onClick={() => handleDeleteClick(bot.id)}
                                    className="text-red-600 hover:underline"
                                >
                                    Delete
                                </button>
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>

            <Pagination page={page} pageSize={pageSize} total={total} onPageChange={handlePageChange} />

            <Modal
                visible={modalVisible}
                title={editingBot ? "Edit Bot" : "Add Bot"}
                onClose={() => setModalVisible(false)}
            >
                <input type="hidden" value={form.id} />
                <div className="mb-4">
                    <input
                        type="text"
                        placeholder="Address"
                        value={form.address}
                        onChange={(e) => setForm({ ...form, address: e.target.value })}
                        className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-400"
                    />
                </div>
                <div className="mb-4">
          <textarea
              placeholder="CA File"
              value={form.ca_file}
              onChange={(e) => setForm({ ...form, ca_file: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-400"
              rows={5}
          />
                </div>
                <div className="mb-4">
          <textarea
              placeholder="KEY File"
              value={form.key_file}
              onChange={(e) => setForm({ ...form, key_file: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-400"
              rows={5}
          />
                </div>
                <div className="mb-4">
          <textarea
              placeholder="CRT File"
              value={form.crt_file}
              onChange={(e) => setForm({ ...form, crt_file: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-400"
              rows={5}
          />
                </div>
                <div className="flex justify-end space-x-2">
                    <button
                        onClick={() => setModalVisible(false)}
                        className="bg-gray-300 hover:bg-gray-400 text-gray-800 px-4 py-2 rounded"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={handleSave}
                        className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded"
                    >
                        Save
                    </button>
                </div>
            </Modal>

            <Modal
                visible={rawConfigVisible}
                title="Command"
                onClose={() => setRawConfigVisible(false)}
            >
        <pre className="max-h-[500px] overflow-y-auto text-sm text-gray-700 whitespace-pre-wrap break-words">
          {rawConfigText.split(/\s+/).filter(Boolean).join("\n")}
        </pre>
            </Modal>

            <Modal
                visible={structuredConfigVisible}
                title="Edit Config"
                onClose={() => setStructuredConfigVisible(false)}
            >
                <ConfigForm botId={selectId} />
            </Modal>

            <ConfirmModal
                visible={confirmVisible}
                message="Are you sure you want to delete this bot?"
                onConfirm={confirmDelete}
                onCancel={cancelDelete}
            />
        </div>
    );
}

export default Bots;
