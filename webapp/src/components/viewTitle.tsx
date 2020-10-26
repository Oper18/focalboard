// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {injectIntl, IntlShape, FormattedMessage} from 'react-intl'

import {BlockIcons} from '../blockIcons'
import {Board} from '../blocks/board'
import mutator from '../mutator'
import Menu from '../widgets/menu'
import MenuWrapper from '../widgets/menuWrapper'
import Editable from '../widgets/editable'

import Button from './button'

type Props = {
    board: Board
    intl: IntlShape
}

type State = {
    isHoverOnCover: boolean
    title: string
}

class ViewTitle extends React.Component<Props, State> {
    private titleEditor = React.createRef<Editable>()
    shouldComponentUpdate(): boolean {
        return true
    }

    constructor(props: Props) {
        super(props)
        this.state = {isHoverOnCover: false, title: props.board.title}
    }

    render(): JSX.Element {
        const {board, intl} = this.props

        return (
            <>
                <div
                    className='octo-hovercontrols'
                    onMouseOver={() => {
                        this.setState({isHoverOnCover: true})
                    }}
                    onMouseLeave={() => {
                        this.setState({isHoverOnCover: false})
                    }}
                >
                    <Button
                        style={{display: (!board.icon && this.state.isHoverOnCover) ? null : 'none'}}
                        onClick={() => {
                            const newIcon = BlockIcons.shared.randomIcon()
                            mutator.changeIcon(board, newIcon)
                        }}
                    >
                        <FormattedMessage
                            id='TableComponent.add-icon'
                            defaultMessage='Add Icon'
                        />
                    </Button>
                </div>

                <div className='octo-icontitle'>
                    {board.icon &&
                        <MenuWrapper>
                            <div className='octo-button octo-icon'>{board.icon}</div>
                            <Menu>
                                <Menu.Text
                                    id='random'
                                    name={intl.formatMessage({id: 'ViewTitle.random-icon', defaultMessage: 'Random'})}
                                    onClick={() => mutator.changeIcon(board, BlockIcons.shared.randomIcon())}
                                />
                                <Menu.Text
                                    id='remove'
                                    name={intl.formatMessage({id: 'ViewTitle.remove-icon', defaultMessage: 'Remove Icon'})}
                                    onClick={() => mutator.changeIcon(board, undefined, 'remove icon')}
                                />
                            </Menu>
                        </MenuWrapper>}
                    <Editable
                        ref={this.titleEditor}
                        className='title'
                        value={this.state.title}
                        placeholderText={intl.formatMessage({id: 'ViewTitle.untitled-board', defaultMessage: 'Untitled Board'})}
                        onChange={(title) => this.setState({title})}
                        onBlur={() => mutator.changeTitle(board, this.state.title)}
                        onFocus={() => this.titleEditor.current.focus()}
                        onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>): void => {
                            if (e.keyCode === 27 && !(e.metaKey || e.ctrlKey) && !e.shiftKey && !e.altKey) { // ESC
                                e.stopPropagation()
                                this.setState({title: this.props.board.title})
                                setTimeout(() => this.titleEditor.current.blur(), 0)
                            } else if (e.keyCode === 13 && !(e.metaKey || e.ctrlKey) && !e.shiftKey && !e.altKey) { // Return
                                e.stopPropagation()
                                mutator.changeTitle(board, this.state.title)
                                this.titleEditor.current.blur()
                            }
                        }}
                    />
                </div>
            </>
        )
    }
}

export default injectIntl(ViewTitle)
