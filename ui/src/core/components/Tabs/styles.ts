/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import styled, { css } from 'styled-components';

type TabItemProps = {
  isActive?: boolean;
};

const TabList = styled.div`
  display: flex;
  border-bottom: 2px solid #79797a;
`;

const TabPanel = styled.div`
  padding: 20px 0;
`;

const TabItem = styled.div<TabItemProps>`
  padding: 0 26px 18px 26px;
  position: relative;
  top: 2px;
  cursor: pointer;

  ${({ isActive }) =>
    isActive &&
    css`
      border-bottom: 2px solid #f4f4fa;
    `}
`;

export default {
  TabList,
  TabPanel,
  TabItem
};
