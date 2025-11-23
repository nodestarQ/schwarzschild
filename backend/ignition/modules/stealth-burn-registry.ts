
import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
//@ts-ignore hardhat ignition does not understand file extensions
import { StealthBurnRegistryContractName } from "../../src/constants";

export default buildModule("stealthBurnRegistry", (m) => {
    const StealthBurnRegistry = m.contract(StealthBurnRegistryContractName, [], { libraries: {} });

    return { StealthBurnRegistry };
});